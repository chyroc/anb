package config

type ConfigTaskUpload struct {
	Src       string `yaml:"src"`
	Dest      string `yaml:"dest"`
	ExpendEnv bool   `yaml:"expend_env"`
}
