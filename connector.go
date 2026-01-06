package shun

import (
	easyssh "github.com/appleboy/easyssh-proxy"
	"golang.org/x/crypto/ssh"
)

// Connector defines methods for establishing SSH connections.
type Connector interface {
	Connect(params ConnectionParams) (*ssh.Client, error)
}

type ConnectionParams struct {
	User string
	Ip   string
	Port string
}

// Connect establishes an SSH connection to the specified user and IP address.
func (r *shun) Connect(params ConnectionParams) (*ssh.Client, error) {
	cfg := easyssh.MakeConfig{
		User:    params.User,
		Server:  params.Ip,
		Port:    params.Port,
		KeyPath: r.KeyFile,
	}

	s, c, err := cfg.Connect()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	r.client = c

	return c, err
}
