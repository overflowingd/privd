package logic

import (
	"net"

	"github.com/google/nftables"
)

type Conn struct {
	*nftables.Conn

	ready        bool
	inet         *nftables.Table
	ip4Whitelist *nftables.Set
}

func New() *Conn {
	return &Conn{
		Conn: &nftables.Conn{},
	}
}

func (c *Conn) Init(inet *nftables.Table, ip4whitel *nftables.Set) {
	if inet == nil || ip4whitel == nil {
		panic("conn: init attempt with nil args")
	}

	c.inet = inet
	c.ip4Whitelist = ip4whitel
	c.ready = true
}

func (c *Conn) Ready() bool {
	return c.ready
}

func (c *Conn) WhitelistIPs(ips ...net.IP) error {
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
		if err := c.SetAddElements(c.ip4Whitelist, elements); err != nil {
			return err
		}
	}

	if len(elements6) > 0 {
		return ErrIp6NotSupported
	}

	return nil
}

func (c *Conn) Flushing(fn func(*Conn) error) error {
	if err := fn(c); err != nil {
		return err
	}

	return c.Flush()
}
