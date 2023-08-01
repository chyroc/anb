package config

func (r *ConfigServer) ServerHost() string {
	return r.User + "@" + r.Host
}

type ConfigServer struct {
	User              string `yaml:"user"`
	Host              string `yaml:"host"`
	SSHPrivateKeyPath string `yaml:"ssh_private_key_path"`
}
