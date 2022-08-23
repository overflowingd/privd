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

	hosts := handler.NewHosts(conn, resolver)
	{
		gr := s.Group("/v1")
		gr.POST("hosts/trusted", hosts.AddTrusted)
	}

	return s
}

func ConfigPath() string {
	if path := os.Getenv(conf.EnvConfigPath); path != "" {
		return path
	}

	return conf.DefaultConfigPath
}

func Config(path string) (*conf.Config, error) {
	cnf := &conf.Config{}

	if err := cleanenv.ReadConfig(path, cnf); err != nil {
		return nil, err
	}

	return cnf, nil
}

func Init(cfg *conf.Config, conn *logic.Conn) error {
	tables, err := conn.ListTables()
	if err != nil {
		return err
	}

	for i := range tables {
		if tables[i].Name == cfg.Table {
			table := tables[i]

			trustedHosts, err := conn.GetSetByName(table, cfg.Sets.TrustedHosts)
			if err != nil {
				return fmt.Errorf("%w: %v", logic.ErrSetNotFound, err)
			}

			trustedHosts6, err := conn.GetSetByName(table, cfg.Sets.TrustedHosts6)
			if err != nil {
				return fmt.Errorf("%w: %v", logic.ErrSetNotFound, err)
			}

			conn.Table = table
			conn.TrustedHosts = trustedHosts
			conn.TrustedHosts6 = trustedHosts6
			return nil
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

	if err := Init(conf, conn); err != nil {
		log.Fatalf("init: %v", err)
	}

	resolverPool := pond.New(conf.Resolver.MaxWorkers, conf.Resolver.MaxWorkers, pond.MinWorkers(conf.Resolver.MinWorkers))

	if err := srv.Start(conf.Srv.Listen, Handler(conn, logic.NewResolver(resolverPool))); err != nil {
		log.Fatalf("srv.start: %v", err)
	}
}
