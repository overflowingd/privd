package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alitto/pond"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"ovfl.io/overflowingd/privd/pkg/conf"
	"ovfl.io/overflowingd/privd/pkg/handler"
	"ovfl.io/overflowingd/privd/pkg/logic"
	"ovfl.io/overflowingd/privd/pkg/srv"
)

func Handler(conn *logic.Conn, resolver *logic.Resolver) *gin.Engine {
	s := gin.Default()

	hosts := handler.NewIp4(conn, resolver)
	{
		gr := s.Group("/v1")
		gr.POST("ip4/whitelists", hosts.Whitelist)
	}

	return s
}

func ConfigPath() string {
	if path := os.Getenv(conf.EnvConfigPath); path != "" {
		return path
	}

	return conf.DefaultPath
}

func Config(path string) (*conf.Config, error) {
	cnf := &conf.Config{}

	if err := cleanenv.ReadConfig(path, cnf); err != nil {
		return nil, err
	}

	cnf.Tables = &conf.Tables{
		Inet: logic.TableInet,
	}

	cnf.Sets = &conf.Sets{
		Ip4Whitelist: logic.Ip4WhitelistSet,
	}

	return cnf, nil
}

func Init(cnf *conf.Config, conn *logic.Conn, resolver *logic.Resolver) error {
	tables, err := conn.ListTables()
	if err != nil {
		return err
	}

	for i := range tables {
		if tables[i].Name == cnf.Tables.Inet {
			table := tables[i]

			ip4Whitelist, err := conn.GetSetByName(table, cnf.Sets.Ip4Whitelist)
			if err != nil {
				return fmt.Errorf("%w: %v", logic.ErrSetNotFound, err)
			}

			conn.TableInet = table
			conn.Ip4WhitelistSet = ip4Whitelist
			return nil
		}
	}

	ips, domains := logic.SplitHosts(cnf.Ip4Whitelist)

	if err := conn.WhitelistIPs(ips...); err != nil {
		if err != logic.ErrIp6NotSupported {
			return err
		}
	}

	if err := conn.Flush(); err != nil {
		return err
	}

	if len(domains) > 0 {
		resolved, err := resolver.Resolve(domains)
		if err != nil {
			return err
		}

		if err := conn.WhitelistIPs(resolved...); err != nil {
			if err != logic.ErrIp6NotSupported {
				return err
			}
		}

		if err := conn.Flush(); err != nil {
			return err
		}
	}

	return logic.ErrTableRequired
}

func main() {
	conf, err := Config(ConfigPath())
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	conn := logic.New()
	resolverPool := pond.New(conf.Resolver.MaxWorkers, conf.Resolver.MaxWorkers, pond.MinWorkers(conf.Resolver.MinWorkers))
	resolver := logic.NewResolver(resolverPool)

	if err := Init(conf, conn, resolver); err != nil {
		log.Fatalf("init: %v", err)
	}

	if err := srv.Start(conf.Srv.Listen, Handler(conn, resolver)); err != nil {
		log.Fatalf("srv.start: %v", err)
	}
}
