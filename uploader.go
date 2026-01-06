package shun

import (
	"bytes"
	"io"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Uploader is responsible for uploading content to a remote private server.
type Uploader interface {

	// Upload writes bytes to a file in dst on a remote server.
	Upload(data []byte, dst string) error
}

type SftpUploader struct {
	sshClient *ssh.Client
}

func NewSftpUploader(s *ssh.Client) *SftpUploader {
	return &SftpUploader{sshClient: s}
}

// Upload writes bytes to a file on a remote server using SFTP.
func (s *SftpUploader) Upload(data []byte, dst string) error {
	sftp, err := sftp.NewClient(s.sshClient)
	if err != nil {
		return err
	}
	defer sftp.Close()

	srcFile := bytes.NewReader(data)
	dstFile, err := sftp.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
