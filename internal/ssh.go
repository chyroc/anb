package internal

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/chyroc/chaos"
	"github.com/hnakamur/go-scp"
	"golang.org/x/crypto/ssh"
)

type SSH struct {
	host string
	user string

	client *ssh.Client

	initAuthMethodOnce sync.Once
	authMethod         ssh.AuthMethod
}

func (r *SSH) Client() *ssh.Client {
	return r.client
}

type SSHConfig struct {
	Host string
	User string
}

func NewSSH(config *SSHConfig) *SSH {
	return &SSH{
		host: config.Host,
		user: config.User,
	}
}

func (r *SSH) Dial() error {
	r.initAuthMethod()

	client, err := ssh.Dial("tcp", r.host+":22", &ssh.ClientConfig{
		User:              r.user,
		Auth:              []ssh.AuthMethod{r.authMethod},
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
		BannerCallback:    nil,
		ClientVersion:     "",
		HostKeyAlgorithms: nil,
		Timeout:           0,
	})
	if err != nil {
		return fmt.Errorf("dial fail: %w", err)
	}

	r.client = client
	return nil
}

func (r *SSH) Close() error {
	return r.client.Close()
}

func (r *SSH) Run(cmd string, args ...interface{}) (string, error) {
	cmd = fmt.Sprintf(cmd, args...)
	session, err := r.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("new session fail: %w", err)
	}

	var bufout bytes.Buffer
	var buferr bytes.Buffer
	session.Stdout = &bufout
	session.Stderr = &buferr
	if err := session.Run(cmd); err != nil {
		ser := strings.TrimSpace(buferr.String())
		if ser != "" {
			return "", fmt.Errorf("run %q fail: %s", cmd, ser)
		}
		return bufout.String(), fmt.Errorf("run %q fail: %w", cmd, err)
	}

	return bufout.String(), nil
}

func (r *SSH) RunInPipe(cmd string, args ...interface{}) (string, error) {
	cmd = fmt.Sprintf(cmd, args...)
	session, err := r.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("new session fail: %w", err)
	}

	bufout := new(bytes.Buffer)
	buferr := new(bytes.Buffer)
	session.Stdout = chaos.TeeWriter([]io.Writer{bufout, os.Stdout}, nil)
	session.Stderr = chaos.TeeWriter([]io.Writer{buferr, os.Stderr}, nil)
	if err := session.Run(cmd); err != nil {
		ser := strings.TrimSpace(buferr.String())
		if ser != "" {
			return "", fmt.Errorf("run %q fail: %s", cmd, ser)
		}
		return bufout.String(), fmt.Errorf("run %q fail: %w", cmd, err)
	}

	return bufout.String(), nil
}

// https://stackoverflow.com/questions/53256373/sending-file-over-ssh-in-go
// https://itectec.com/unixlinux/ssh-the-protocol-for-sending-files-over-ssh-in-code/
func (r *SSH) WriteFile(bs []byte, filemode string, filename string) (finalErr error) {
	session, err := r.client.NewSession()
	if err != nil {
		return fmt.Errorf("new session fail: %w", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	dir, base := filepath.Split(filename)

	go func() {
		defer wg.Done()

		stdin, err := session.StdinPipe()
		if err != nil {
			finalErr = err
			return
		}
		defer stdin.Close()

		if _, err = fmt.Fprintf(stdin, "C%s %d %s\n", filemode, len(bs), base); err != nil {
			finalErr = err
			return
		}
		if _, err = stdin.Write(bs); err != nil {
			finalErr = err
			return
		}
		if _, err = fmt.Fprint(stdin, "\x00"); err != nil {
			finalErr = err
			return
		}
	}()

	if finalErr != nil {
		return finalErr
	}

	if err = session.Run("/usr/bin/scp -qt " + dir); err != nil {
		return err
	}
	wg.Wait()

	return nil
}

func (r *SSH) initAuthMethod() {
	r.initAuthMethodOnce.Do(func() {
		var signers []ssh.Signer
		dir := os.Getenv("HOME") + "/.ssh/"
		fs, _ := ioutil.ReadDir(dir)
		for _, f := range fs {
			if strings.HasSuffix(f.Name(), ".pub") {
				continue
			}
			if pubInfo, _ := os.Stat(dir + f.Name() + ".pub"); pubInfo == nil {
				continue
			}
			data, err := ioutil.ReadFile(dir + f.Name())
			if err != nil {
				continue
			}
			signer, err := ssh.ParsePrivateKey(data)
			if err != nil {
				continue
			}
			signers = append(signers, signer)
		}
		r.authMethod = ssh.PublicKeys(signers...)
	})
}

func (r *SSH) Upload(src, dest string) error {
	stat, err := os.Stat(src) // scp 不支持 ln，所以这里暂时这么写
	if err != nil {
		return err
	}

	rr := scp.NewSCP(r.client)
	if stat.IsDir() {
		return rr.SendDir(src, dest, func(parentDir string, info os.FileInfo) (bool, error) {
			if info.IsDir() {
				return true, nil
			}
			localPath := parentDir + "/" + info.Name()
			remotePath := GetRemoteRevPath(src, dest, localPath, false)
			if r.isServerFileMd5Equal(localPath, remotePath) {
				PrintfGreen("\t[upload] %q skip\n", src)
				return false, nil
			} else {
				PrintfYellow("\t[upload] %q running...\n", src)
				return true, nil
			}
		})
	} else {
		if r.isServerFileMd5Equal(src, dest) {
			PrintfGreen("\t[upload] %q skip\n", src)
			return nil
		} else {
			PrintfYellow("\t[upload] %q running...\n", src)
			return rr.SendFile(src, dest)
		}
	}
}

func (r *SSH) Download(src, dest string) error {
	out, err := r.Run("ls -ld %q | awk '{print $1}'", src)
	if err != nil {
		return err
	}
	isDir := strings.TrimSpace(out) != "" && strings.TrimSpace(out)[0] == 'd'

	rr := scp.NewSCP(r.client)
	if isDir {
		return rr.ReceiveDir(src, dest, func(parentDir string, info os.FileInfo) (bool, error) {
			if info.IsDir() {
				return true, nil
			}
			localPath := parentDir + "/" + info.Name()
			remotePath := GetRemoteRevPath(dest, src, localPath, true)
			if r.isServerFileMd5Equal(localPath, remotePath) {
				PrintfGreen("\t[download] %q skip\n", src)
				return false, nil
			} else {
				PrintfYellow("\t[download] %q running...\n", src)
				return true, nil
			}
		})
	} else {
		if r.isServerFileMd5Equal(dest, src) {
			PrintfGreen("\t[download] %q skip\n", src)
			return nil
		} else {
			PrintfYellow("\t[download] %q running...\n", src)
			return rr.ReceiveFile(src, dest)
		}
	}
}

func (r *SSH) PrintMeta() {
	fmt.Printf("--- ssh meta ---\n")
	fmt.Printf("user: %s\n", r.client.User())
	fmt.Printf("session: %x\n", r.client.SessionID())
	fmt.Printf("client-version: %s\n", r.client.ClientVersion())
	fmt.Printf("server-version: %s\n", r.client.ServerVersion())
	fmt.Printf("remove-addr: %s\n", r.client.RemoteAddr())
	fmt.Printf("local-addr: %s\n", r.client.LocalAddr())
	fmt.Printf("--- ssh meta ---\n\n")
}

func (r *SSH) isServerFileMd5Equal(local, remote string) bool {
	localMd5, _ := GetFileMd5(local)
	sshMd5, _ := r.sshGetFileMd5(remote)
	return sshMd5 != "" && localMd5 == sshMd5
}

func (r *SSH) sshGetFileMd5(file string) (string, error) {
	out, err := r.Run("md5sum %s", file)
	if err != nil {
		return "", err
	}
	ss := strings.Split(out, " ")
	if len(ss) >= 2 {
		return ss[0], nil
	}
	return "", fmt.Errorf("invalid md5: %q", out)
}
