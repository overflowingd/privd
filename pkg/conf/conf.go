package conf

const (
	DefaultPath = ".privd.yml"
)

const (
	EnvConfigPath = "CNF_PATH"
)

type Config struct {
	Srv          *Srv      `yaml:"srv"`
	Resolver     *Resolver `yaml:"resolver"`
	Nft          *Nft      `yaml:"nft"`
	Ip4Whitelist []string  `yaml:"ip4-whitel"`
	Tables       *Tables
	Sets         *Sets
}

type Nft struct {
	RulesUp   string `yaml:"rules-up"`
	RulesDown string `yaml:"rules-down"`
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
