// rshell defines methods for remote shell operations.
// It currently abstracts SSH connection and file upload functionalities.
package shun

import (
	"golang.org/x/crypto/ssh"
)

type Shun interface {
	Connector
	Uploader
}

type shun struct {
	KeyFile string
	client  *ssh.Client
}

// NewShun creates a new instance of rshell with the provided SSH key file.
func NewShun(keyFile string) Shun {
	return &shun{
		KeyFile: keyFile,
	}
}
