package logic

import (
	"net"

	"github.com/google/nftables"
)

type Conn struct {
	*nftables.Conn

	Table        *nftables.Table
	TrustedHosts *nftables.Set
}

func New() *Conn {
	return &Conn{
		Conn: &nftables.Conn{},
	}
}

func (c *Conn) IP(addr string) net.IP {
	ip := net.ParseIP(addr)
	if ip == nil {
		return nil
	}

	if ip.To4() == nil {
		return ip.To16()
	}

	return ip.To4()
}
