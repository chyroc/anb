package config

func (r *ConfigServer) ServerHost() string {
	return r.User + "@" + r.Host
}

type ConfigServer struct {
	User string `yaml:"user"`
	Host string `yaml:"host"`
}
