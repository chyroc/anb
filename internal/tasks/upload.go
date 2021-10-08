package tasks

import (
	"io/ioutil"
	"os"
	"regexp"

	"github.com/chyroc/anb/internal"
	"github.com/chyroc/anb/internal/config"
	"github.com/google/uuid"
)

func RunUploadTask(task *config.Task, cli *internal.SSHCommand) error {
	d := task.Upload
	internal.PrintfWhite("\t[upload] %q => %q\n", d.Src, d.Dest)
	src := d.Src
	if d.ExpendEnv {
		src = expendFileEnv(src)
	}
	if err := cli.UploadAnyFile(src, d.Dest); err != nil {
		return err
	}
	return nil
}

func expendFileEnv(file string) string {
	dir, _ := ioutil.TempDir("", "anb-tmp-*")
	if dir == "" {
		return file
	}
	bs, _ := ioutil.ReadFile(file)
	if len(bs) == 0 {
		return file
	}

	filename := dir + "/" + uuid.New().String()
	state, _ := os.Lstat(file)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, state.Mode())
	if err != nil {
		return file
	}
	_, _ = f.WriteString(os.ExpandEnv(string(bs)))
	f.Close()

	return filename
}

func splitWithColon(s string) (string, string) {
	m := regSplitWitholon.FindStringSubmatch(s)
	if len(m) == 3 {
		return m[1], m[2]
	}
	return "", ""
}

var regSplitWitholon = regexp.MustCompile(`^(.*?)\s*:\s*(.*?)$`)
