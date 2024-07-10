package dns

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/SergeyCherepiuk/dns-go/internal/dns/serde"
	"github.com/SergeyCherepiuk/dns-go/internal/dns/types"
)

func StartServer(ctx context.Context, addr net.UDPAddr) error {
	cache := NewDNSCache(ctx)

	for {
		conn, err := net.ListenUDP("udp", &addr)
		if err != nil {
			return err
		}

		err = handleConnection(conn, cache)
		if err != nil {
			return err
		}
	}
}

func handleConnection(conn *net.UDPConn, cache *dnsCache) error {
	defer conn.Close()

	buf := make([]byte, 512)
	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}

	query, err := serde.UnmarshalPacket(buf[:n])
	if err != nil {
		return err
	}

	fmt.Println(query.String())

	response, ok := lookupCache(query, cache)
	if !ok {
		response, err = Lookup(query)
		if err != nil {
			return err
		}

		var (
			domain  = response.Answers[0].Domain
			answers = response.Answers
			ttl     = time.Duration(response.Answers[0].Ttl) * time.Second
		)
		cache.set(domain, answers, ttl)
	}

	fmt.Println(response.String())

	responseBytes, err := serde.MarshalPacket(response)
	if err != nil {
		return err
	}

	_, err = conn.WriteToUDP(responseBytes, addr)
	return err
}

func lookupCache(query types.Packet, cache *dnsCache) (types.Packet, bool) {
	domain := query.Questions[0].Domain
	record, ok := cache.get(domain)
	if !ok {
		return types.Packet{}, false
	}

	response := types.Packet{
		Header: types.Header{
			ID:                  query.Header.ID,
			PacketType:          types.PacketTypeResponse,
			RecursionDesired:    query.Header.RecursionDesired,
			RecursionAvailable:  true,
			QuestionSectionSize: query.Header.QuestionSectionSize,
			AnswerSectionSize:   uint16(len(record.answers)),
		},
		Questions: query.Questions,
		Answers:   record.answers,
	}

	return response, true
}
