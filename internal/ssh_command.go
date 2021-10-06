package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type SSHCommand struct {
	Ins *SSH
}

func NewSSHCommand(sshIns *SSH) *SSHCommand {
	return &SSHCommand{
		Ins: sshIns,
	}
}

func (r *SSHCommand) Dial() error {
	return r.Ins.Dial()
}

func (r *SSHCommand) Run(cmd string) (string, error) {
	return r.Ins.Run(cmd)
}

func (r *SSHCommand) Close() error {
	return r.Ins.Close()
}

func (r *SSHCommand) Md5File(file string) (string, error) {
	out, err := r.Ins.Run("md5sum " + file)
	if err != nil {
		return "", err
	}
	ss := strings.Split(out, " ")
	if len(ss) >= 2 {
		return ss[0], nil
	}
	return "", fmt.Errorf("md5sum response %q not valid", out)
}

func (r *SSHCommand) CreateDir(dir, filemode string) error {
	_, err := r.Ins.Run(fmt.Sprintf("mkdir -m %s -p %s", filemode, dir))
	return err
}

func (r *SSHCommand) CopyFile(src, dest string) (bool, error) {
	f, _ := os.Stat(src)

	destMd5, _ := r.Md5File(dest)
	srcMd5, _ := NewLocalCommand().Md5File(src)
	if srcMd5 != "" && srcMd5 == destMd5 {
		return false, nil
	}
	bs, err := ioutil.ReadFile(src)
	if err != nil {
		return false, err
	}

	if err = r.Ins.WriteFile(bs, GetFilePerm(f.Mode()), dest); err != nil {
		return false, err
	}
	return true, nil
}
