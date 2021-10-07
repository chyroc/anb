package config

import (
	"fmt"
	"strings"
)

type Task struct {
	// 公用 task 参数
	ID   string `yaml:"id"`   // 唯一性 ID，不允许重复，可以通过 $tasks.<id>.outputs.<name> 取数据
	Name string `yaml:"name"` // 名称，打印用
	If   IfExpr `yaml:"if"`   // 一个 expr 表达式，为真，则这个 task 才可以运行
	Dir  string `yaml:"dir"`  // 执行 command 的目录

	// 下面是 task 的具体实现
	Copy     *ConfigTaskCopy `yaml:"copy"`
	Cmd      *ConfigTaskCmd  `yaml:"cmd"`
	LocalCmd *ConfigTaskCmd  `yaml:"local_cmd"`
}

// 是否可以运行
func (r *Task) ShouldRun(vals Any) (string, bool, error) {
	if r.If == "" {
		return "", true, nil
	}

	isTrue, err := r.If.IsTrue(vals)
	if err != nil {
		return "", false, fmt.Errorf("%q invalid: %w", r.If, err)
	}
	return strings.TrimSpace(string(r.If)), isTrue, nil
}

// 任务名称
func (r *Task) TaskName(idx int) string {
	if r.Name != "" {
		return r.Name
	}
	return fmt.Sprintf("task #%d", idx+1)
}
