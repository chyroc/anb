package internal

import (
	"fmt"
)

type RunRequest struct {
	Config string
}

func Run(req *RunRequest) error {
	conf, err := LoadConfig(req.Config)
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
	PrintfWhite("[server] %s connected\n", conf.ServerHost())

	for _, task := range conf.Tasks {
		switch task.TaskType() {
		case TaskTypeCopy:
			d := task.Copy
			PrintfWhite("[copy] %q => %q\n", d.Src, d.Dest)
			return cli.CopyAnyFile(d.Src, d.Dest)
		case TaskTypeCmd:
			for _, cmd := range task.Cmd.Commands {
				PrintfWhite("[cmd] %q\n", cmd)
				out, err := cli.Run(cmd)
				if err != nil {
					return err
				}
				fmt.Print(out)
			}
		default:
			return fmt.Errorf("不支持 " + string(task.TaskType()))
		}
	}
	return nil
}
