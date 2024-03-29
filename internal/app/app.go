package app

import (
	"fmt"
	"io/ioutil"

	"github.com/chyroc/anb/internal"
	"github.com/chyroc/anb/internal/config"
	"github.com/chyroc/anb/internal/tasks"
)

type RunRequest struct {
	Config string
}

func Run(req *RunRequest) error {
	conf, err := config.LoadConfig(req.Config)
	if err != nil {
		return err
	}

	var privateKey []byte
	if conf.Server.SSHPrivateKeyPath != "" {
		if privateKey, err = ioutil.ReadFile(conf.Server.SSHPrivateKeyPath); err != nil {
			return err
		}
	}

	cli := internal.NewSSH(&internal.SSHConfig{
		User:          conf.Server.User,
		Host:          conf.Server.Host,
		SSHPrivateKey: privateKey,
	})
	defer cli.Close()

	if err = cli.Dial(); err != nil {
		return err
	}
	internal.PrintfWhite("[server] %s connected\n", conf.Server.ServerHost())

	vals := config.Any{"$tasks": config.Any{}}

	for idx, task := range conf.Tasks {
		internal.PrintfWhite("[task] %s\n", task.TaskName(idx))

		if reason, shouldRun, err := task.ShouldRun(vals); err != nil {
			return err
		} else if !shouldRun {
			internal.PrintfGreen("\t[skip] %s\n", reason)
			continue
		}
		switch task.TaskType() {
		case config.TaskTypeUpload:
			if err := tasks.RunUploadTask(task, cli); err != nil {
				return err
			}
		case config.TaskTypeDownload:
			if err := tasks.RunDownloadTask(task, cli, vals); err != nil {
				return err
			}
		case config.TaskTypeCmd:
			if err := tasks.RunCmd(task, cli, vals); err != nil {
				return err
			}
		case config.TaskTypeLocalCmd:
			if err := tasks.RunLocalCmd(task, cli, vals); err != nil {
				return err
			}
		default:
			return fmt.Errorf("不支持命令: " + string(task.TaskType()))
		}
	}
	return nil
}
