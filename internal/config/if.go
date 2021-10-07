package config

import (
	"os"
	"strings"

	"github.com/chyroc/anbko/env"
	"github.com/chyroc/anbko/vm"
)

// IfExpr 是 if 表达式，输入: string，输出: bool
//
// IfExpr 支持下面几个函数和操作符
// 	函数：exist
// 	操作符：&& || !
//
// 示例
// if: exist("/path")
// if: exist(/path)
// if: !exist(/path)
// if: exist(/path1) && !exist(/path2)
// if: (exist(/path1) || !exist(/path2)) && exist(/path3)
type IfExpr string

func (r IfExpr) IsTrue(vals Any) (bool, error) {
	s := strings.TrimSpace(string(r))
	return ExecBoolScript(s, vals)
}

type Any map[string]interface{}

func ExecBoolScript(script string, val Any) (bool, error) {
	e := env.NewEnv()

	if err := e.Define("exist", func(path string) bool {
		f, _ := os.Lstat(path)
		return f != nil
	}); err != nil {
		return false, err
	}
	// err := e.DefineGlobal("$ab", "/Users/chyroc/code/chyroc/anb/README.md")
	for k, v := range val {
		if err := e.DefineGlobal(k, v); err != nil {
			return false, err
		}
	}
	res, err := vm.Execute(e, nil, script)
	if err != nil {
		return false, err
	}
	b, _ := res.(bool)
	return b, nil
}
