package dns

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"

	"github.com/SergeyCherepiuk/dns-go/internal/dns/serde"
	"github.com/SergeyCherepiuk/dns-go/internal/dns/types"
)

var RootServers = [13]net.IP{
	net.IPv4(198, 41, 0, 4),     // a.root-servers.net.
	net.IPv4(170, 247, 170, 2),  // b.root-servers.net.
	net.IPv4(192, 33, 4, 12),    // c.root-servers.net.
	net.IPv4(199, 7, 91, 13),    // d.root-servers.net.
	net.IPv4(192, 203, 230, 10), // e.root-servers.net.
	net.IPv4(192, 5, 5, 241),    // f.root-servers.net.
	net.IPv4(192, 112, 36, 4),   // g.root-servers.net.
	net.IPv4(198, 97, 190, 53),  // h.root-servers.net.
	net.IPv4(192, 36, 148, 17),  // i.root-servers.net.
	net.IPv4(192, 58, 128, 30),  // j.root-servers.net.
	net.IPv4(193, 0, 14, 129),   // k.root-servers.net.
	net.IPv4(199, 7, 83, 42),    // l.root-servers.net.
	net.IPv4(202, 12, 27, 33),   // m.root-servers.net.
}

var (
	ErrUnableToResolve   = errors.New("unable to resolve")
	ErrInvalidRecordType = errors.New("invalid record type")
)

func Lookup(query types.Packet) (types.Packet, error) {
	var (
		initialDomain   = query.Questions[0].Domain
		rootServerIndex = rand.Intn(len(RootServers))
		rootServerIP    = RootServers[rootServerIndex]
		addr            = net.UDPAddr{IP: rootServerIP, Port: 53}
	)

	for {
		response, err := sendQuery(query, addr)
		if err != nil {
			return types.Packet{}, err
		}

		if response.Header.ResponseCode != types.ResponseCodeNoError {
			return response, nil
		}

		if _, ok := getIPv4(response.Answers, initialDomain); ok {
			return response, nil
		}

		cname, ok := getCname(response.Answers, initialDomain)
		if ok {
			cnameQuery := constructQuery(cname)
			cnameResponse, err := Lookup(cnameQuery)
			if err != nil {
				return types.Packet{}, err
			}

			response.Answers = append(response.Answers, cnameResponse.Answers...)
			response.Header.AnswerSectionSize += cnameResponse.Header.AnswerSectionSize
			return response, nil
		}

		host, ok := pickNameServer(response.AuthorityRecords)
		if !ok {
			return types.Packet{}, ErrUnableToResolve
		}

		addr.IP, ok = resolveNameServer(response.AdditionalRecords, host)
		if ok {
			continue
		}

		nsQuery := constructQuery(host)
		hostResponse, err := Lookup(nsQuery)
		if err != nil {
			return types.Packet{}, err
		}

		addr.IP, ok = getIPv4(hostResponse.Answers, host)
		if !ok {
			return types.Packet{}, ErrUnableToResolve
		}
	}
}

func sendQuery(query types.Packet, addr net.UDPAddr) (types.Packet, error) {
	localAddr := net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0}

	conn, err := net.DialUDP("udp", &localAddr, &addr)
	if err != nil {
		return types.Packet{}, err
	}

	queryBytes, err := serde.MarshalPacket(query)
	if err != nil {
		return types.Packet{}, err
	}

	n, err := conn.Write(queryBytes)
	if err != nil {
		return types.Packet{}, err
	}

	if n != len(queryBytes) {
		err = fmt.Errorf("unread bytes (server read %d out of %d)", n, len(queryBytes))
		return types.Packet{}, err
	}

	responseBytes := make([]byte, 512)
	n, err = conn.Read(responseBytes)
	if err != nil {
		return types.Packet{}, err
	}

	return serde.UnmarshalPacket(responseBytes)
}

func getIPv4(records []types.Record, domain string) (net.IP, bool) {
	for _, record := range records {
		if record.Type == types.RecordTypeA && record.Domain == domain {
			return record.Data.(net.IP), true
		}
	}
	return nil, false
}

func getCname(records []types.Record, domain string) (string, bool) {
	for _, record := range records {
		if record.Type == types.RecordTypeCNAME && record.Domain == domain {
			return record.Data.(string), true
		}
	}
	return "", false
}

func pickNameServer(records []types.Record) (string, bool) {
	for _, record := range records {
		if record.Type == types.RecordTypeNS {
			host := record.Data.(string)
			return host, true
		}
	}
	return "", false
}

func resolveNameServer(records []types.Record, host string) (net.IP, bool) {
	for _, record := range records {
		if record.Type == types.RecordTypeA && record.Domain == host {
			ip := record.Data.(net.IP)
			return ip, true
		}
	}
	return nil, false
}

func constructQuery(domain string) types.Packet {
	return types.Packet{
		Header: types.Header{
			ID:                  uint16(rand.Intn(math.MaxUint16)),
			PacketType:          types.PacketTypeQuery,
			Opcode:              types.OpcodeQuery,
			QuestionSectionSize: 1,
		},
		Questions: []types.Question{
			{
				Domain: domain,
				Type:   types.QuestionTypeA,
				Class:  types.QuestionClassIN,
			},
		},
	}
}
