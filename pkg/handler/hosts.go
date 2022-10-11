package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ovfl.io/overflowingd/privd/pkg/logic"
)

type Ip4 struct {
	conn     *logic.Conn
	resolver *logic.Resolver
}

func NewIp4(conn *logic.Conn, resolver *logic.Resolver) *Ip4 {
	return &Ip4{
		conn:     conn,
		resolver: resolver,
	}
}

func (r *Ip4) Whitelist(ctx *gin.Context) {
	items := make([]string, 0)

	if err := ctx.BindJSON(&items); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ips, domains := logic.SplitHosts(items)

	resolved, err := r.resolver.Resolve(domains)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ips = append(ips, resolved...)

	{
		err := r.conn.Flushing(func(c *logic.Conn) error {
			err := c.WhitelistIPs(ips...)
			if err != logic.ErrIp6NotSupported {
				return err
			}

			return nil
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
	}
}
