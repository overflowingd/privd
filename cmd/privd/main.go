package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

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

func ApplyNftRules(path string) error {
	cmd := exec.Command("nft", "-f", path)

	stderr := bytes.NewBuffer(make([]byte, 0))
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}

	return nil
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

			conn.Init(table, ip4Whitelist)
			break
		}
	}

	if !conn.Ready() {
		return logic.ErrTableRequired
	}

	ips, domains := logic.SplitHosts(cnf.Ip4Whitelist)

	err = conn.Flushing(func(c *logic.Conn) error {
		return c.WhitelistIPs(ips...)
	})

	if err != nil {
		if err != logic.ErrIp6NotSupported {
			return err
		}
	}

	if len(domains) > 0 {
		resolved, err := resolver.Resolve(domains)
		if err != nil {
			return err
		}

		err = conn.Flushing(func(c *logic.Conn) error {
			return c.WhitelistIPs(resolved...)
		})

		if err != nil {
			if err != logic.ErrIp6NotSupported {
				return err
			}
		}
	}

	return nil
}

func main() {
	cnf, err := Config(ConfigPath())
	if err != nil {
		log.Printf("config: %v", err)
		return
	}

	conn := logic.New()
	resolverPool := pond.New(cnf.Resolver.MaxWorkers, cnf.Resolver.MaxWorkers, pond.MinWorkers(cnf.Resolver.MinWorkers))
	resolver := logic.NewResolver(resolverPool)

	if err := ApplyNftRules(cnf.Nft.RulesUp); err != nil {
		log.Printf("load-nft-rules: %v", err)
		return
	}

	defer ApplyNftRules(cnf.Nft.RulesDown)

	if err := Init(cnf, conn, resolver); err != nil {
		log.Printf("init: %v", err)
		return
	}

	if err := srv.Start(cnf.Srv.Listen, Handler(conn, resolver)); err != nil {
		log.Printf("srv.start: %v", err)
		return
	}
}
