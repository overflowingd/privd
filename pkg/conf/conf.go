package conf

const (
	DefaultPath = ".cnf.yml"
)

const (
	EnvConfigPath = "CNF_PATH"
)

type Config struct {
	Tables       *Tables   `yaml:"tables"`
	Sets         *Sets     `yaml:"sets"`
	Srv          *Srv      `yaml:"srv"`
	Resolver     *Resolver `yaml:"resolver"`
	Ip4Whitelist []string  `yaml:"ip4_whitel"`
}

type Tables struct {
	Inet string
}

type Sets struct {
	Ip4Whitelist string
}

type Srv struct {
	Listen string `yaml:"listen" env:"PRIVD_LISTEN"`
}

type Resolver struct {
	MaxWorkers int `yaml:"max-workers" env:"PRIVD_RESOLV_MAX_WORKERS" env-default:"10"`
	MinWorkers int `yaml:"min-workers" env:"PRIVD_RESOLV_MIN_WORKERS" env-default:"1"`
}
