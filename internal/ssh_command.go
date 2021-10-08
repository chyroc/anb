package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
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

func (r *SSHCommand) CopyAnyFile(src, dest string) error {
	f, _ := os.Lstat(src)
	if f.IsDir() {
		if err := r.CreateDir(dest, GetFilePerm(f.Mode())); err != nil {
			return err
		}
		return Walk(src, dest, func(isDir bool, path, target string) error {
			return r.CopyAnyFile(path, target)
		})
	} else if f.Mode()&os.ModeSymlink != 0 {
		link, err := os.Readlink(src)
		if err != nil {
			return err
		}
		created, err := r.CreateHardLink(link, dest)
		if err != nil {
			return err
		}
		if created {
			PrintfYellow("\t[copy] %q running...\n", src)
		} else {
			PrintfGreen("\t[copy] %q skip\n", src)
		}
	} else {
		fd, _ := os.Lstat(filepath.Dir(src))
		if err := r.CreateDir(filepath.Dir(dest), GetFilePerm(fd.Mode())); err != nil {
			return err
		}
		copied, err := r.CopyFile(src, dest)
		if err != nil {
			return err
		}
		if copied {
			PrintfYellow("\t[copy] %q running...\n", src)
		} else {
			PrintfGreen("\t[copy] %q skip\n", src)
		}
	}
	return nil
}

func (r *SSHCommand) MoveFile(src, dest string) error {
	_, err := r.Ins.Run("mv %s %s", src, dest)
	return err
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

	fakeDest := "/tmp/" + uuid.New().String()
	if err = r.Ins.WriteFile(bs, GetFilePerm(f.Mode()), fakeDest); err != nil {
		return false, err
	}

	if err := r.MoveFile(fakeDest, dest); err != nil {
		return false, err
	}

	// if err = r.Ins.WriteFile(bs, GetFilePerm(f.Mode()), dest); err != nil {
	// 	return false, err
	// }
	return true, nil
}

func (r *SSHCommand) CreateHardLink(source, target string) (bool, error) {
	// a -> b
	out, err := r.Ins.Run("ls -l " + target)
	if out != "" {
		if strings.Contains(out, " -> "+source) {
			return false, nil
		}
		return false, fmt.Errorf("%q exist", target)
	}
	_, err = r.Ins.Run(fmt.Sprintf("ln -s %s %s", source, target))
	return true, err
}
