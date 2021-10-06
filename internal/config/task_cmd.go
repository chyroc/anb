package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type ConfigTaskCmd struct {
	Commands []string
}

func (r *ConfigTaskCmd) UnmarshalYAML(unmarshal func(interface{}) error) error {
	{
		var data []string
		if err := unmarshal(&data); err == nil {
			r.Commands = data
			return nil
		}
	}
	{
		var data string
		if err := unmarshal(&data); err == nil {
			r.Commands = append(r.Commands, data)
			return nil
		}
	}
	return fmt.Errorf("不支持的 cmd 命令")
}

func (r ConfigTaskCmd) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(r.Commands)
}
