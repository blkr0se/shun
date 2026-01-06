// shun defines operations for interacting with remote private servers via ssh.
package shun

import (
	"os"

	"golang.org/x/crypto/ssh"
)

// Shun is the core operator for interacting with remote private servers.
// By default, it uses SSH to connect to a server and SFTP to upload content.
type Shun struct {
	// Connector is used as a factory for SSH connections, on which
	// many of Shun's services have a dependency on.
	Connector *SshConnector

	// Cached SSH client, set during the Connect() method call.
	sshClient *ssh.Client

	// Uploader handles file upload operations to a remote destination server.
	Uploader Uploader
}

// Upload copies a local file from src to a remote dst.
func (shun *Shun) Upload(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return shun.Uploader.Upload(data, dst)
}

// Connect establishes an SSH connection to a remote server.
func (shun *Shun) Connect() error {
	client, err := shun.Connector.Connect()
	if err != nil {
		return err
	}

	shun.sshClient = client

	return nil
}
