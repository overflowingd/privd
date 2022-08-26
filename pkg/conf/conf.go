package conf

const (
	DefaultPath = "cfg.yml"
)

const (
	EnvConfigPath = "CFG_PATH"
)

const (
	Table = "bd1ae1fc24d1d1e9c3e864a79f176c56"
)

type Config struct {
	Table    string    `yaml:"table" env:"PRIVD_NFTABLE"`
	Sets     *Sets     `yaml:"sets"`
	Srv      *Srv      `yaml:"srv"`
	Resolver *Resolver `yaml:"resolver"`
}

type Sets struct {
	Ns              string
	TrustedHosts    string
	TrustedHosts6   string
	NontunneledNets string
}

type Srv struct {
	Listen string `yaml:"listen" env:"PRIVD_LISTEN"`
}

type Resolver struct {
	MaxWorkers int `yaml:"max-workers" env:"PRIVD_RESOLV_MAX_WORKERS" env-default:"10"`
	MinWorkers int `yaml:"min-workers" env:"PRIVD_RESOLV_MIN_WORKERS" env-default:"1"`
}
