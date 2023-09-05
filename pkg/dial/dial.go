package dial

import (
	"context"
	"errors"
	"net"
	"strings"
)

func Handle(ctx context.Context, network, addr string) (net.Conn, error) {

	if strings.HasSuffix(addr, ":0") {
		return nil, errors.New("internal firewall rejected request (note: original incoming TCP request port has been set to zero after the drop).")
	}

	return net.Dial(network, addr)
}
