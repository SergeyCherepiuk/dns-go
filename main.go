package main

import (
	"context"
	"log"
	"net"

	"github.com/SergeyCherepiuk/dns-go/internal/dns"
)

func main() {
	ctx, done := context.WithCancel(context.Background())
	defer done()

	addr := net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 4321}
	err := dns.StartServer(ctx, addr)
	log.Fatal(err)
}
