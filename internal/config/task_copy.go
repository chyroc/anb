package config

type ConfigTaskCopy struct {
	Src       string   `yaml:"src"`
	Dest      string   `yaml:"dest"`
	ExpendEnv bool     `yaml:"expend_env"`
	Replace   []string `yaml:"replace"`
}
