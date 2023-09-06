package rewriter

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"

	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/statute"
	"scaleship.io/kubernetes/kubeception-proxy/pkg/encryption"
)

type ClusterDynamicRewriter struct {
	Private *rsa.PrivateKey
}

func (c ClusterDynamicRewriter) Rewrite(ctx context.Context, r *socks5.Request) (context.Context, *statute.AddrSpec) {

	// Drop direct IP usage
	if r.DstAddr.FQDN == "" {
		r.DestAddr.Port = 0
		return ctx, r.DestAddr
	}

	// Force API service port to ":443"
	r.DestAddr.Port = 443

	endpoint := fmt.Sprintf("%s:%d", r.DestAddr.FQDN, r.DestAddr.Port)

	// Verify authentication signature with endpoint
	signed, err := encryption.Verify(&c.Private.PublicKey, endpoint, r.AuthContext.Payload["signature"])
	if !signed || err != nil {
		log.Printf("[WARNING] Invalid signatures from %q in SOCKS5 authentication process.", r.RemoteAddr)
		r.DestAddr.Port = 0
		return ctx, r.DestAddr
	}

	return ctx, r.DestAddr
}
