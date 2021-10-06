package internal

import (
	"fmt"
	"os"
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

			srcInfo, err := os.Stat(d.Src)
			if err != nil {
				return err
			}
			if srcInfo.IsDir() {
				err = Walk(d.Src, d.Dest, func(isDir bool, path, target string) error {
					f, _ := os.Stat(path)
					if isDir {
						return cli.CreateDir(target, GetFilePerm(f.Mode()))
					}
					copied, err := cli.CopyFile(path, target)
					if err != nil {
						return err
					}
					if copied {
						PrintfYellow("\tcopy: %q\n", path)
					} else {
						PrintfGreen("\tskip: %q\n", path)
					}
					return nil
				})
				return err
			} else {
				copied, err := cli.CopyFile(d.Src, d.Dest)
				if err != nil {
					return err
				}
				if copied {
					PrintfYellow("\tcopy: %q\n", d.Src)
				} else {
					PrintfGreen("\tskip: %q\n", d.Src)
				}
			}
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
