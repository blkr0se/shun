package shun

import (
	easyssh "github.com/appleboy/easyssh-proxy"
	"golang.org/x/crypto/ssh"
)

type SshConnector struct {
	User    string
	Ip      string
	Port    string
	KeyFile string
}

// Connect establishes an SSH connection to the specified user and IP address.
func (c *SshConnector) Connect() (*ssh.Client, error) {
	cfg := easyssh.MakeConfig{
		User:    c.User,
		Server:  c.Ip,
		Port:    c.Port,
		KeyPath: c.KeyFile,
	}

	s, conn, err := cfg.Connect()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	return conn, nil
}
