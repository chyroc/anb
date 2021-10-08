package tasks

import (
	"github.com/chyroc/anb/internal"
	"github.com/chyroc/anb/internal/config"
)

func RunDownloadTask(task *config.Task, cli *internal.SSH, vals config.Any) error {
	d := task.Download
	internal.PrintfWhite("\t[download] %q => %q\n", d.Src, d.Dest)
	return cli.Download(d.Src, d.Dest)
}
