package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server *ConfigServer `yaml:"server"`
	Tasks  []*Task       `yaml:"tasks"`
}

func LoadConfig(file string) (*Config, error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	conf := new(Config)
	err = yaml.Unmarshal(bs, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
