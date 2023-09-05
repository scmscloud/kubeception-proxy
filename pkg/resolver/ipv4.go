package resolver

import (
	"context"
	"errors"
	"fmt"
	"net"
)

type IPv4 struct{}

func (c IPv4) Resolve(ctx context.Context, addr string) (context.Context, net.IP, error) {

	ips, err := net.LookupIP(addr)
	if err != nil {
		return ctx, nil, err
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			return ctx, ip, nil
		}
	}

	return ctx, nil, errors.New(fmt.Sprintf("Unable to resolve address %q.", addr))
}
