package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

type SSH struct {
	host string
	user string

	client *ssh.Client

	initAuthMethodOnce sync.Once
	authMethod         ssh.AuthMethod
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

func (r *SSH) Run(cmd string) (string, error) {
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

// https://stackoverflow.com/questions/53256373/sending-file-over-ssh-in-go
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
