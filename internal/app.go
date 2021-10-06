package internal

import (
	"fmt"
	"io/ioutil"
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
				return fmt.Errorf("不支持 文件夹")
			} else {
				destMd5, err := cli.Md5File(d.Dest)
				if err != nil {
					return err
				}
				srcMd5, err := NewLocalCommand().Md5File(d.Src)
				if err != nil {
					return err
				}
				if srcMd5 == destMd5 {
					PrintfGreen("\tequal file, skip\n")
					continue
				}
				bs, err := ioutil.ReadFile(d.Src)
				if err != nil {
					return err
				}
				err = cli.Ins.WriteFile(bs, GetFilePerm(srcInfo.Mode()), d.Dest)
				if err != nil {
					return err
				}
				PrintfYellow("\tcopy...\n")
			}
		case TaskTypeCmd:
			for _, cmd := range task.Cmd.Commands {
				PrintfWhite("[cmd] %q\n", cmd)
				out, err := cli.Run(cmd)
				if err != nil {
					return err
				}
				fmt.Println(out)
			}
		default:
			return fmt.Errorf("不支持 " + string(task.TaskType()))
		}
	}
	return nil
}
