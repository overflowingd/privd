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

func (r *Resolver) Resolve(hosts []string) (ips []net.IP, err error) {
	ips = make([]net.IP, 0, len(hosts))

	group := r.pool.Group()

	for _, h := range hosts {
		group.Submit(func() {
			a, e := net.LookupHost(h)
			if e != nil {
				err = e
				return
			}

			for i := range a {
				ips = append(ips, NewIP(a[i]))
			}
		})
	}

	group.Wait()

	return
}

func NewIP(addr string) net.IP {
	ip := net.ParseIP(addr)
	if ip == nil {
		return nil
	}

	if ip.To4() == nil {
		return ip.To16()
	}

	return ip.To4()
}

func SplitHosts(hosts []string) (ips []net.IP, domains []string) {
	for _, host := range hosts {
		if ip := NewIP(host); ip != nil {
			ips = append(ips, ip)
			continue
		}

		domains = append(domains, host)
	}

	return
}
