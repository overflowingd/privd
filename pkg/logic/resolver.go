package logic

import (
	"net"

	"github.com/alitto/pond"
)

type Resolver struct {
	pool *pond.WorkerPool
}

func NewResolver(p *pond.WorkerPool) *Resolver {
	return &Resolver{
		pool: p,
	}
}

func (r *Resolver) Resolve(hosts []string) (addrs []string, err error) {
	addrs = make([]string, 0, len(hosts))

	group := r.pool.Group()

	for _, h := range hosts {
		group.Submit(func() {
			a, e := net.LookupHost(h)
			if e != nil {
				err = e
				return
			}

			addrs = append(addrs, a...)
		})
	}

	group.Wait()
	return
}
