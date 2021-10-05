package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
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

func (r *SSH) Run(cmd string) (string, error) {
	session, err := r.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("new session fail: %w", err)
	}

	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	if err := session.Run(cmd); err != nil {
		return "", fmt.Errorf("run %q fail: %w", cmd, err)
	}

	return buf.String(), nil
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
