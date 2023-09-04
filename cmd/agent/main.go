package main

import (
	"flag"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/net/proxy"
)

var (
	addr   = flag.String("bind-address", "127.0.0.1:8080", "The IP:PORT address on which to listen.")
	remote = flag.String("proxy-address", "127.0.0.1:8000", "Address (IP:PORT) of the kubeception remote server.")
)

func main() {

	log.SetFlags(0)
	flag.Parse()

	log.Printf("Proxy is now listening on %q...", *addr)
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Cannot listen: %s", err.Error())
	}

	for {

		conn, err := lis.Accept()
		if err != nil {
			log.Printf("[WARNING] Cannot accept: %s", err.Error())
		}

		go func(src net.Conn) {

			defer src.Close()

			dailer, err := proxy.SOCKS5("tcp", *remote, nil, &net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 30 * time.Second,
			})

			if err != nil {
				log.Printf("[WARNING] Cannot initialize %q SOCKS5 proxy: %s", *remote, err.Error())
				return
			}

			// Advertise tcp address is always defined by remote proxy (`localhost:443` is overwritten).
			dst, err := dailer.Dial("tcp", "localhost:443")
			if err != nil {
				log.Printf("[WARNING] Cannot dial: %s", err.Error())
				log.Println(err)
				return
			}

			done := make(chan struct{})

			go forward(src, dst, done)
			go forward(dst, src, done)

			<-done
			<-done
		}(conn)
	}
}

func forward(src net.Conn, dst net.Conn, done chan<- struct{}) {

	defer src.Close()
	defer dst.Close()

	io.Copy(src, dst)

	done <- struct{}{}

}
