package dns

import (
	"context"
	"fmt"
	"net"

	"github.com/SergeyCherepiuk/dns-go/internal/dns/cache"
	"github.com/SergeyCherepiuk/dns-go/internal/dns/serde"
)

func StartServer(ctx context.Context, addr net.UDPAddr) error {
	cache := cache.NewDnsCache(ctx)

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

func handleConnection(conn *net.UDPConn, cache *cache.DnsCache) error {
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

	response, err := Lookup(query, cache)
	if err != nil {
		return err
	}

	fmt.Println(response.String())

	responseBytes, err := serde.MarshalPacket(response)
	if err != nil {
		return err
	}

	_, err = conn.WriteToUDP(responseBytes, addr)
	return err
}
