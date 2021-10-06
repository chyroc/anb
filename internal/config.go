package internal

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

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

type Config struct {
	Server *ConfigServer `yaml:"server"`
	Tasks  []*ConfigTask `yaml:"tasks"`
}

func (r *Config) ServerHost() string {
	return r.Server.User + "@" + r.Server.Host
}

type ConfigServer struct {
	User string `yaml:"user"`
	Host string `yaml:"host"`
}

type ConfigTask struct {
	Name string          `yaml:"name"`
	Copy *ConfigTaskCopy `yaml:"copy"`
	Cmd  *ConfigTaskCmd  `yaml:"cmd"`
}

type ConfigTaskCopy struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}

type ConfigTaskCmd struct {
	Commands []string
}

type TaskType string

const (
	TaskTypeCopy TaskType = "copy"
	TaskTypeCmd  TaskType = "cmd"
)

func (r *ConfigTask) TaskType() TaskType {
	if r.Copy != nil {
		return TaskTypeCopy
	}
	if r.Cmd != nil {
		return TaskTypeCmd
	}
	return ""
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
