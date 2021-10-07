package internal

import (
	"fmt"

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

	cli := NewSSHCommand(NewSSH(&SSHConfig{
		User: conf.Server.User,
		Host: conf.Server.Host,
	}))
	defer cli.Close()

	if err = cli.Dial(); err != nil {
		return err
	}
	PrintfWhite("[server] %s connected\n", conf.Server.ServerHost())

	vals := config.Any{"$tasks": config.Any{}}

	for idx, task := range conf.Tasks {
		PrintfWhite("[task] %s\n", task.TaskName(idx))

		if reason, shouldRun, err := task.ShouldRun(vals); err != nil {
			return err
		} else if !shouldRun {
			PrintfGreen("\t[skip] %s\n", reason)
			continue
		}
		switch task.TaskType() {
		case config.TaskTypeCopy:
			if err := tasks.RunCopyTask(task, cli); err != nil {
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
			return fmt.Errorf("不支持 " + string(task.TaskType()))
		}
	}
	return nil
}
