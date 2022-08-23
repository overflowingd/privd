package conf

const (
	DefaultConfigPath = ".cfg.yml"
)

const (
	EnvConfigPath = "CFG_PATH"
)

type Config struct {
	Table    string    `yaml:"table" env:"PRIVD_NFTABLE"`
	Sets     *Sets     `yaml:"sets"`
	Srv      *Srv      `yaml:"srv"`
	Resolver *Resolver `yaml:"resolver"`
}

type Sets struct {
	TrustedHosts  string `yaml:"trusted-hosts" env:"PRIVD_TRUSTED_HOSTS"`
	TrustedHosts6 string `yaml:"trusted-hosts6" env:"PRIVD_TRUSTED_HOSTS6"`
}

type Srv struct {
	Listen string `yaml:"listen" env:"PRIVD_LISTEN"`
}

type Resolver struct {
	MaxWorkers int `yaml:"max-workers" env:"PRIVD_RESOLV_MAX_WORKERS" env-default:"10"`
	MinWorkers int `yaml:"min-workers" env:"PRIVD_RESOLV_MIN_WORKERS" env-default:"1"`
}
