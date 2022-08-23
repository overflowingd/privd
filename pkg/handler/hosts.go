package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/nftables"
	"ovfl.io/overflowingd/privd/pkg/logic"
)

type Hosts struct {
	conn     *logic.Conn
	resolver *logic.Resolver
}

func NewHosts(conn *logic.Conn, resolver *logic.Resolver) *Hosts {
	return &Hosts{
		conn:     conn,
		resolver: resolver,
	}
}

func (h *Hosts) AddTrusted(ctx *gin.Context) {
	trusted := make([]string, 0)

	if err := ctx.BindJSON(&trusted); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	hosts := make([]string, 0, len(trusted))
	elements := make([]nftables.SetElement, 0, len(trusted))

	for _, tr := range trusted {
		if ip := h.conn.IP(tr); ip != nil {
			elements = append(elements, nftables.SetElement{Key: []byte(ip)})
			continue
		}

		hosts = append(hosts, tr)
	}

	addrs, err := h.resolver.Resolve(hosts)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, addr := range addrs {
		// todo:v1: handle ipv6
		elements = append(elements, nftables.SetElement{Key: []byte(h.conn.IP(addr))})
	}

	if err := h.conn.SetAddElements(h.conn.TrustedHosts, elements); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.conn.Flush(); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
