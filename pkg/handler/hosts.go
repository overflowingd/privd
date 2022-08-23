package handler

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
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
	items := make([]string, 0)

	if err := ctx.BindJSON(&items); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	hosts := make([]string, 0, len(items))
	ips := make([]net.IP, 0, len(items))

	for _, item := range items {
		if ip := h.conn.IP(item); ip != nil {
			ips = append(ips, ip)
			continue
		}

		hosts = append(hosts, item)
	}

	addrs, err := h.resolver.Resolve(hosts)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	for _, addr := range addrs {
		ips = append(ips, h.conn.IP(addr))
	}

	if err := h.conn.AddTrustedHosts(ips...); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if err := h.conn.Flush(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
}
