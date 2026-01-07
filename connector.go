package shun

import (
	"os"
	"os/exec"

	easyssh "github.com/appleboy/easyssh-proxy"
	"golang.org/x/crypto/ssh"
)

type SshConnector struct {
	User        string
	Ip          string
	Port        string
	KeyFile     string
	UseTerminal bool
}

// Connect establishes an SSH connection to the specified user and IP address.
func (c *SshConnector) Connect() (*ssh.Client, error) {
	cfg := easyssh.MakeConfig{
		User:       c.User,
		Server:     c.Ip,
		Port:       c.Port,
		KeyPath:    c.KeyFile,
		RequestPty: c.UseTerminal,
	}

	s, conn, err := cfg.Connect()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	return conn, nil
}

// Interactive starts a remote interactive shell using OpenSSH's client.
func (c *SshConnector) Interactive(args ...string) error {
	sshCmd := exec.Command("ssh", args...)

	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	if err := sshCmd.Start(); err != nil {
		return err
	}

	if err := sshCmd.Wait(); err != nil {
		return err
	}

	return nil
}
