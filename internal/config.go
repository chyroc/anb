package internal

import (
	"fmt"
	"io/ioutil"
	"os"

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
	Name       string          `yaml:"name"`
	IfNotExist string          `yaml:"if_not_exist"`
	IfExist    string          `yaml:"if_exist"`
	Copy       *ConfigTaskCopy `yaml:"copy"`
	Cmd        *ConfigTaskCmd  `yaml:"cmd"`
	LocalCmd   *ConfigTaskCmd  `yaml:"local_cmd"`
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
	TaskTypeCopy     TaskType = "copy"
	TaskTypeCmd      TaskType = "cmd"
	TaskTypeLocalCmd TaskType = "local_cmd"
)

func (r *ConfigTask) TaskType() TaskType {
	if r.Copy != nil {
		return TaskTypeCopy
	}
	if r.Cmd != nil {
		return TaskTypeCmd
	}
	if r.LocalCmd != nil {
		return TaskTypeLocalCmd
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

func (r *ConfigTask) ShouldRun() (string, bool) {
	if r.IfNotExist != "" {
		// 只有不存在，则运行（true）
		if f, _ := os.Lstat(r.IfNotExist); f != nil {
			return fmt.Sprintf("%q should not exist, skip", r.IfNotExist), false
		}
	}

	if r.IfExist != "" {
		if f, _ := os.Lstat(r.IfNotExist); f == nil {
			return fmt.Sprintf("%q should exist, skip", r.IfExist), false
		}
	}

	return "", true
}

func (r *ConfigTask) TaskName(idx int) string {
	if r.Name != "" {
		return r.Name
	}
	return fmt.Sprintf("task #%d", idx+1)
}
