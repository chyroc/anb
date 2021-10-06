package internal

import (
	"fmt"
	"strings"
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

	for idx, task := range conf.Tasks {
		PrintfWhite("[task] %s\n", task.TaskName(idx))

		if reason, shouldRun := task.ShouldRun(); !shouldRun {
			PrintfGreen("\t%s\n", reason)
			continue
		}
		switch task.TaskType() {
		case TaskTypeCopy:
			d := task.Copy
			PrintfWhite("\t[copy] %q => %q\n", d.Src, d.Dest)
			if err := cli.CopyAnyFile(d.Src, d.Dest); err != nil {
				return err
			}
		case TaskTypeCmd:
			for _, cmd := range task.Cmd.Commands {
				PrintfWhite("\t[cmd] %q\n", cmd)
				out, err := cli.Run(cmd)
				if err != nil {
					return err
				}
				fmt.Print(out)
			}
		case TaskTypeLocalCmd:
			for _, cmd := range task.LocalCmd.Commands {
				cmds := strings.Fields(cmd)
				PrintfWhite("\t[local_cmd] %q\n", cmd)
				fmt.Println(cmds)
				out, err := NewLocalCommand().RunCommand(cmds[0], cmds[1:]...)
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
