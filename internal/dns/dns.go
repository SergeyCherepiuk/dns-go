package dns

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"
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

func Lookup(query Packet) (Packet, error) {
	var (
		initialDomain   = query.Questions[0].Domain
		rootServerIndex = rand.Intn(len(RootServers))
		rootServerIP    = RootServers[rootServerIndex]
		addr            = net.UDPAddr{IP: rootServerIP, Port: 53}
	)

	for {
		response, err := sendQuery(query, addr)
		if err != nil {
			return Packet{}, err
		}

		if _, ok := getIPv4(response); ok {
			return response, nil
		}

		cname, ok := getCanonicalName(response, initialDomain)
		if ok {
			cnameQuery := constructCnameQuery(cname)
			cnameResponse, err := Lookup(cnameQuery)
			if err != nil {
				return Packet{}, err
			}

			response.Answers = append(response.Answers, cnameResponse.Answers...)
			response.Header.AnswerSectionSize += cnameResponse.Header.AnswerSectionSize
			return response, nil
		}

		host, err := pickNameServer(response)
		if err != nil {
			return Packet{}, err
		}

		addr.IP, err = resolveNameServer(host, response)
		if err != nil {
			return Packet{}, err
		}
	}
}

func sendQuery(query Packet, addr net.UDPAddr) (Packet, error) {
	localAddr := net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0}

	conn, err := net.DialUDP("udp", &localAddr, &addr)
	if err != nil {
		return Packet{}, err
	}

	queryBytes, err := MarshalPacket(query)
	if err != nil {
		return Packet{}, err
	}

	n, err := conn.Write(queryBytes)
	if err != nil {
		return Packet{}, err
	}

	if n != len(queryBytes) {
		err = fmt.Errorf("unread bytes (server read %d out of %d)", n, len(queryBytes))
		return Packet{}, err
	}

	responseBytes := make([]byte, 512)
	n, err = conn.Read(responseBytes)
	if err != nil {
		return Packet{}, err
	}

	return UnmarshalPacket(responseBytes)
}

func getCanonicalName(response Packet, domain string) (string, bool) {
	for _, answer := range response.Answers {
		if answer.Type == RecordTypeCNAME && answer.Domain == domain {
			return answer.Data.(string), true
		}
	}
	return "", false
}

func getIPv4(response Packet) (net.IP, bool) {
	for _, answer := range response.Answers {
		if answer.Type == RecordTypeA {
			return answer.Data.(net.IP), true
		}
	}
	return nil, false
}

func pickNameServer(response Packet) (string, error) {
	for _, authorityRecord := range response.AuthorityRecords {
		if authorityRecord.Type == RecordTypeNS {
			host := authorityRecord.Data.(string)
			return host, nil
		}
	}
	return "", ErrUnableToResolve
}

func resolveNameServer(host string, response Packet) (net.IP, error) {
	for _, additionalRecord := range response.AdditionalRecords {
		if additionalRecord.Type == RecordTypeA && additionalRecord.Domain == host {
			ip := additionalRecord.Data.(net.IP)
			return ip, nil
		}
	}

	hostQuery := Packet{
		Header: Header{
			ID:                  uint16(rand.Intn(math.MaxUint16)),
			PacketType:          PacketTypeQuery,
			Opcode:              OpcodeQuery,
			RecursionDesired:    true,
			QuestionSectionSize: 1,
		},
		Questions: []Question{
			{
				Domain: host,
				Type:   QuestionTypeA,
				Class:  QuestionClassIN,
			},
		},
	}

	hostResponse, err := Lookup(hostQuery)
	if err != nil {
		return nil, err
	}

	ip, ok := getIPv4(hostResponse)
	if !ok {
		return nil, ErrUnableToResolve
	}

	return ip, nil
}

func constructCnameQuery(cname string) Packet {
	return Packet{
		Header: Header{
			ID:                  uint16(rand.Intn(math.MaxUint16)),
			PacketType:          PacketTypeQuery,
			Opcode:              OpcodeQuery,
			QuestionSectionSize: 1,
		},
		Questions: []Question{
			{
				Domain: cname,
				Type:   QuestionTypeA,
				Class:  QuestionClassIN,
			},
		},
	}
}
