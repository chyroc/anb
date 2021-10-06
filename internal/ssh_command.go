package internal

import (
	"fmt"
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
