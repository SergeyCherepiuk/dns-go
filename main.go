package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SergeyCherepiuk/dns-go/internal/dns"
)

func main() {
	addr := net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 4321}
	err := listen(addr)
	log.Fatal(err)
}

func listen(addr net.UDPAddr) error {
	for {
		conn, err := net.ListenUDP("udp", &addr)
		if err != nil {
			return err
		}

		err = handleConnection(conn)
		if err != nil {
			return err
		}
	}
}

func handleConnection(conn *net.UDPConn) error {
	defer conn.Close()

	buf := make([]byte, 512)
	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}

	query, err := dns.UnmarshalPacket(buf[:n])
	if err != nil {
		return err
	}

	fmt.Println(query.String())

	response, err := lookup(query)
	if err != nil {
		return err
	}

	fmt.Println(response.String())

	responseBytes, err := dns.MarshalPacket(response)
	if err != nil {
		return err
	}

	_, err = conn.WriteToUDP(responseBytes, addr)
	return err
}

func lookup(query dns.Packet) (dns.Packet, error) {
	var (
		listenAddr  = net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0}
		receiveAddr = net.UDPAddr{IP: net.IPv4(8, 8, 8, 8), Port: 53}
	)

	conn, err := net.DialUDP("udp", &listenAddr, &receiveAddr)
	if err != nil {
		return dns.Packet{}, err
	}

	queryBytes, err := dns.MarshalPacket(query)
	if err != nil {
		return dns.Packet{}, err
	}

	n, err := conn.Write(queryBytes)
	if err != nil {
		return dns.Packet{}, err
	}

	if n != len(queryBytes) {
		err = fmt.Errorf("unread bytes (server read %d out of %d)", n, len(queryBytes))
		return dns.Packet{}, err
	}

	responseBytes := make([]byte, 512)
	n, err = conn.Read(responseBytes)
	if err != nil {
		return dns.Packet{}, err
	}

	return dns.UnmarshalPacket(responseBytes)
}
