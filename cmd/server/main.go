package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/things-go/go-socks5"
)

var (
	addr      = flag.String("bind-address", "0.0.0.0:8000", "The IP:PORT address on which to listen.")
	hostname  = flag.String("hostname", "localhost", "Advertised hostname of kubernetes API (must be a DNS `A` record)")
	port      = flag.Int("port", 443, "Advertised port of kubernetes API.")
	refresh   = flag.String("internal-refresh", "5m", "")
	advertise = fmt.Sprintf("%s:%d", *hostname, *port)
)

func main() {

	log.SetFlags(0)
	flag.Parse()

	interval, err := time.ParseDuration(*refresh)
	if err != nil {
		log.Fatalf("Unable to parse refresh duration: %s.", err.Error())
	}

	go func() {
		for {

			current := advertise

			ips, err := net.LookupIP(*hostname)
			if err != nil {
				log.Fatalf("Could not get IPs: %v\n", err)
			}

			for _, ip := range ips {
				if ip.To4() != nil {
					advertise = fmt.Sprintf("%s:%d", ip.String(), *port)
					if advertise != current {
						log.Printf("Advertised IP address of %s as been updated (from %q to %q).", *hostname, current, advertise)
					}
				}
			}

			time.Sleep(interval)
		}
	}()

	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "SOCKS5: ", 0))),
		socks5.WithDial(func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial(network, advertise)
		}),
	)

	log.Printf("Proxy is now listening on %q...", *addr)
	if err := server.ListenAndServe("tcp", *addr); err != nil {
		log.Fatal(err)
	}

}
