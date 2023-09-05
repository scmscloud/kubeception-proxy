package main

import (
	"flag"
	"log"
	"os"

	"github.com/things-go/go-socks5"
	"scaleship.io/kubernetes/kubeception-proxy/pkg/authentication"
	"scaleship.io/kubernetes/kubeception-proxy/pkg/dial"
	"scaleship.io/kubernetes/kubeception-proxy/pkg/encryption"
	"scaleship.io/kubernetes/kubeception-proxy/pkg/resolver"
	"scaleship.io/kubernetes/kubeception-proxy/pkg/rewriter"
)

var (
	addr        = flag.String("bind-address", "0.0.0.0:1080", "The IP:PORT address on which to listen.")
	hostname    = flag.String("hostname", "localhost", "Advertised hostname of kubernetes API (must be a DNS `A` record)")
	privatefile = flag.String("private-file", "", "Path to the private key file for SOCKS5 authentication.")
)

func main() {

	log.SetFlags(0)
	flag.Parse()

	private, err := encryption.ParsePrivate(*privatefile)
	if err != nil {
		log.Fatal(err)
	}

	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "SOCKS5: ", 0))),
		socks5.WithAuthMethods([]socks5.Authenticator{authentication.Signature{}}),
		socks5.WithResolver(resolver.IPv4{}),
		socks5.WithRewriter(rewriter.ClusterDynamicRewriter{
			Private: private,
		}),
		socks5.WithDial(dial.Handle),
	)

	log.Printf("Proxy is now listening on %q...", *addr)
	if err := server.ListenAndServe("tcp", *addr); err != nil {
		log.Fatal(err)
	}

}
