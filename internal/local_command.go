package internal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/chyroc/chaos"
)

type Command interface {
	Md5File(file string) (string, error)
}

type LocalCommand struct{}

func NewLocalCommand() *LocalCommand {
	return &LocalCommand{}
}

func (r *LocalCommand) Md5File(file string) (string, error) {
	out, err := r.RunCommand("md5sum " + file)
	if err != nil {
		return "", err
	}
	ss := strings.Split(out, " ")
	if len(ss) >= 2 {
		return ss[0], nil
	}
	return "", fmt.Errorf("md5sum response %q not valid", out)
}

func (r *LocalCommand) RunCommand(cmd string) (string, error) {
	var bufout bytes.Buffer
	var buferr bytes.Buffer
	c := exec.Command("/bin/sh", "-c", cmd)
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

func (r *LocalCommand) RunCommandInPipe(cmd string) (string, error) {
	bufout := new(bytes.Buffer)
	buferr := new(bytes.Buffer)
	c := exec.Command("/bin/sh", "-c", cmd)
	c.Stdout = chaos.TeeWriter([]io.Writer{bufout, os.Stdout}, nil)
	c.Stderr = chaos.TeeWriter([]io.Writer{buferr, os.Stderr}, nil)
	if err := c.Run(); err != nil {
		ser := strings.TrimSpace(buferr.String())
		if ser != "" {
			return "", fmt.Errorf("run %q fail: %s", cmd, ser)
		}
		return bufout.String(), fmt.Errorf("run %q fail: %w", cmd, err)
	}

	return bufout.String(), nil
}
