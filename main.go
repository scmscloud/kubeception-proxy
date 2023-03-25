package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"log"
	"net"
	"os"
)

var certpool *x509.CertPool
var (
	ip         = flag.String("bind-address", "", "The IP:PORT address on which to listen.")
	advertise  = flag.String("advertise-address", "", "The Kubernetes API Server (\"hostname/ip\":\"port\") public address.")
	cacertfile = flag.String("cacert-file", "", "Path to a client certificate authority cert file for TLS")
	certfile   = flag.String("cert-file", "", "Path to a client cert file for TLS.")
	keyfile    = flag.String("key-file", "", "Path to a client key file for TLS.")
)

func main() {

	flag.Parse()

	cacert, err := os.ReadFile(*cacertfile)
	if err != nil {
		log.Fatal("Unable to read certificate authority: ", err.Error())
	}

	certpool = x509.NewCertPool()
	if ok := certpool.AppendCertsFromPEM(cacert); !ok {
		log.Fatalf("Unable to parse cert from %s", *cacertfile)
	}

	cert, err := os.ReadFile(*certfile)
	if err != nil {
		log.Fatal("Unable to read certificate: ", err.Error())
	}

	key, err := os.ReadFile(*keyfile)

	if err != nil {
		log.Fatal("Unable to read private key: ", err.Error())
	}

	certificate, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{certificate}}

	ln, err := tls.Listen("tcp", *ip, config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Proxy is now listening on \"%s\".", *ip)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go HandleAdvertiseConnection(conn)
	}

}

func HandleAdvertiseConnection(src net.Conn) {

	log.Printf("New connection %v connected via secure channel.", src.RemoteAddr())

	config := &tls.Config{RootCAs: certpool}
	dst, err := tls.Dial("tcp", *advertise, config)
	if err != nil {
		log.Printf("Unable to initiate connection with advertise server: \"%s\".", err.Error())
	}

	done := make(chan struct{})

	go func() {
		defer src.Close()
		defer dst.Close()
		io.Copy(dst, src)
		done <- struct{}{}
	}()

	go func() {
		defer src.Close()
		defer dst.Close()
		io.Copy(src, dst)
		done <- struct{}{}
	}()

	<-done
	<-done

	log.Printf("Connection from %v closed.", src.RemoteAddr())
}
