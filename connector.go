package shun

import (
	"os"
	"os/exec"

	easyssh "github.com/appleboy/easyssh-proxy"
	"golang.org/x/crypto/ssh"
)

type Connector interface {
	Connect() error
	Close() error
	Wait() error
}

type NativeSshProvider struct {
	User    string
	Ip      string
	Port    string
	KeyFile string

	client *ssh.Client
}

// Connect establishes an SSH connection using the native Golang SSH package (golang.org/x/crypto/ssh).
func (nc *NativeSshProvider) Connect() error {
	cfg := easyssh.MakeConfig{
		User:    nc.User,
		Server:  nc.Ip,
		Port:    nc.Port,
		KeyPath: nc.KeyFile,
	}

	s, client, err := cfg.Connect()
	if err != nil {
		return err
	}
	defer s.Close()

	nc.client = client

	return nil
}

func (nc *NativeSshProvider) Close() error {
	return nc.client.Close()
}

func (nc *NativeSshProvider) Wait() error {
	return nc.client.Wait()
}

type BinSshProvider struct {
	// Arguments for the SSH command.
	args []string

	// Command that manages the underlying ssh process.
	cmd *exec.Cmd
}

func NewBinProvider(args ...string) Connector {
	return &BinSshProvider{
		args: args,
	}
}

// Connect establishes a SSH connection using a local OpenSSH binary.
func (bc *BinSshProvider) Connect() error {
	sshCmd := exec.Command("ssh", bc.args...)

	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	bc.cmd = sshCmd

	if err := sshCmd.Start(); err != nil {
		return err
	}

	return nil
}

func (bc *BinSshProvider) Close() error {
	return bc.cmd.Process.Kill()
}

func (bc *BinSshProvider) Wait() error {
	return bc.cmd.Process.Kill()
}

// Deprecated: use NativeProvider instead.
type SshConnector struct {
	User    string
	Ip      string
	Port    string
	KeyFile string
}

// Deprecated: use NativeProvider instead.
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
