package tasks

import (
	"github.com/chyroc/anb/internal"
	"github.com/chyroc/anb/internal/config"
)

func RunCopyTask(task *config.Task, cli *internal.SSHCommand) error {
	d := task.Copy
	internal.PrintfWhite("\t[copy] %q => %q\n", d.Src, d.Dest)
	if err := cli.CopyAnyFile(d.Src, d.Dest); err != nil {
		return err
	}
	return nil
}
