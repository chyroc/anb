package internal

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Command interface {
	Md5File(file string) (string, error)
}

type LocalCommand struct{}

func NewLocalCommand() *LocalCommand {
	return &LocalCommand{}
}

func (r *LocalCommand) Md5File(file string) (string, error) {
	out, err := r.RunCommand("md5sum", file)
	if err != nil {
		return "", err
	}
	ss := strings.Split(out, " ")
	if len(ss) >= 2 {
		return ss[0], nil
	}
	return "", fmt.Errorf("md5sum response %q not valid", out)
}

func (r *LocalCommand) RunCommand(cmd string, args ...string) (string, error) {
	var bufout bytes.Buffer
	var buferr bytes.Buffer
	strings.Split(cmd, " ")
	c := exec.Command(cmd, args...)
	c.Stdout = &bufout
	c.Stderr = &buferr
	if err := c.Run(); err != nil {
		ser := strings.TrimSpace(buferr.String())
		if ser != "" {
			return "", fmt.Errorf("run %q fail: %s", cmd, ser)
		}
		return bufout.String(), fmt.Errorf("run %q fail: %w", cmd, err)
	}

	return bufout.String(), nil
}
