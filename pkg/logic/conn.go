package logic

import (
	"net"

	"github.com/google/nftables"
)

type Conn struct {
	*nftables.Conn

	Table         *nftables.Table
	TrustedHosts  *nftables.Set
	TrustedHosts6 *nftables.Set
}

func New() *Conn {
	return &Conn{
		Conn: &nftables.Conn{},
	}
}

func (c *Conn) AddTrustedHosts(ips ...net.IP) error {
	var ipsLen = len(ips)
	if ipsLen < 1 {
		return nil
	}

	var (
		elements  = make([]nftables.SetElement, 0, ipsLen)
		elements6 = make([]nftables.SetElement, 0, ipsLen/4)
	)

	for _, ip := range ips {
		if ip.To4() != nil {
			elements = append(elements, nftables.SetElement{Key: []byte(ip.To4())})
			continue
		}

		elements6 = append(elements6, nftables.SetElement{Key: []byte(ip.To16())})
	}

	if len(elements) > 0 {
		if err := c.SetAddElements(c.TrustedHosts, elements); err != nil {
			return err
		}
	}

	if len(elements6) > 0 {
		if err := c.SetAddElements(c.TrustedHosts6, elements6); err != nil {
			return err
		}
	}

	return nil
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
