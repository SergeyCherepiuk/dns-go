package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SergeyCherepiuk/dns-go/internal/dns"
)

func main() {
	var (
		listenAddr  = net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0}
		receiveAddr = net.UDPAddr{IP: net.IPv4(8, 8, 8, 8), Port: 53}
	)

	conn, err := net.DialUDP("udp", &listenAddr, &receiveAddr)

	if err != nil {
		log.Fatal(err)
	}

	query := dns.Packet{
		Header: dns.Header{
			ID:                  0x1234,
			PacketType:          dns.PacketTypeQuery,
			Opcode:              dns.OpcodeQuery,
			RecursionDesired:    true,
			QuestionSectionSize: 1,
		},
		Questions: []dns.Question{
			{
				Domain: "google.com.",
				Type:   dns.QuestionTypeA,
				Class:  dns.QuestionClassIN,
			},
		},
	}

	queryBytes := dns.MarshalPacket(query)
	n, err := conn.Write(queryBytes)

	if err != nil {
		log.Fatal(err)
	}

	if n != len(queryBytes) {
		log.Fatalf("not all bytes have been read by the server (%d out of %d)", n, len(queryBytes))
	}

	responseBytes := make([]byte, 512)
	n, err = conn.Read(responseBytes)

	if err != nil {
		log.Fatal(err)
	}

	response := dns.UnmarshalPacket(responseBytes[:n])

	for _, answer := range response.Answers {
		fmt.Printf(
			"%s -> %v (ttl: %d)\n",
			answer.Domain,
			net.IP(answer.Data),
			answer.Ttl,
		)
	}
}
