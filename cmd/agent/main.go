package main

import (
	"flag"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/net/proxy"
	"scaleship.io/kubernetes/kubeception-proxy/pkg/encryption"
)

var (
	addr        = flag.String("bind-address", "127.0.0.1:8080", "The IP:PORT address on which to listen.")
	remote      = flag.String("proxy-address", "127.0.0.1:1080", "Address (IP:PORT) of the kubeception remote server.")
	endpoint    = flag.String("endpoint", "kubernetes.default.svc.cluster.local:443", "Address (IP:PORT) of the kubeception API server.")
	privatefile = flag.String("private-file", "", "Path to the private key file for SOCKS5 authentication.")
)

func main() {

	log.SetFlags(0)
	flag.Parse()

	private, err := encryption.ParsePrivate(*privatefile)
	if err != nil {
		log.Fatal(err)
	}

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

			pwd, err := encryption.Sign(private, *endpoint)
			if err != nil {
				log.Fatal(err)
			}

			gzp, err := encryption.Compress(pwd)

			data := *gzp
			username := data[0 : (len(data)-1)/2]
			password := data[(len(data) / 2):]

			dailer, err := proxy.SOCKS5("tcp", *remote, &proxy.Auth{
				User:     username,
				Password: password,
			}, &net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 30 * time.Second,
			})

			if err != nil {
				log.Printf("[WARNING] Cannot initialize %q SOCKS5 proxy: %s", *remote, err.Error())
				return
			}

			dst, err := dailer.Dial("tcp", *endpoint)
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
