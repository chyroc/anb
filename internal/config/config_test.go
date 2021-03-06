package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Config(t *testing.T) {
	as := assert.New(t)
	chTestDir()

	t.Run("", func(t *testing.T) {
		conf, err := LoadConfig("testdata/config.1.yaml")
		as.Nil(err)
		as.Equal("root", conf.Server.User)
		as.Equal("1.2.3.4", conf.Server.Host)

		{
			as.Equal("x", conf.Tasks[0].Name)
			as.Equal("a", conf.Tasks[0].Upload.Src)
			as.Equal("b", conf.Tasks[0].Upload.Dest)
			as.Equal(TaskTypeUpload, conf.Tasks[0].TaskType())
		}
		{
			as.Equal([]string{"ls", "ls -alh"}, conf.Tasks[1].Cmd.Commands)
			as.Equal(TaskTypeCmd, conf.Tasks[1].TaskType())
		}
		{
			as.Equal([]string{"ls"}, conf.Tasks[2].Cmd.Commands)
			as.Equal(TaskTypeCmd, conf.Tasks[2].TaskType())
		}
	})
}

func chTestDir() {
	pwd, _ := os.Getwd()
	if strings.Contains(pwd, "/internal") {
		count := strings.Count(strings.SplitN(pwd, "/internal", 2)[1], "/")
		for i := 0; i <= count; i++ {
			_ = os.Chdir("..")
		}
	}
}
