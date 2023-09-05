package authentication

import (
	"net"
)

type Client struct {
	Source      string
	Destination Destination
}

type Destination struct {
	FQDN string
	IPv4 net.IP
}
