package shun

import (
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Uploader is responsible for uploading content to a remote private server.
type Uploader interface {

	// Upload writes bytes to a file in dst on a remote server.
	Upload(src string, dst string) error
}

type SftpUploader struct {
	sshClient *ssh.Client

	uploadStrategy RemoteFileCopier
	showProgress   bool
}

func NewSftpUploader(s *ssh.Client) *SftpUploader {
	return &SftpUploader{sshClient: s, uploadStrategy: &AtOnce{}}
}

func (su *SftpUploader) UseProgressBar(v bool) error {
	su.showProgress = v
	su.uploadStrategy = &Chunked{}

	return nil
}

// Upload writes bytes to a file on a remote server using SFTP.
func (su *SftpUploader) Upload(src string, dst string) error {
	client, err := sftp.NewClient(
		su.sshClient,
		sftp.UseConcurrentReads(true),
		sftp.UseConcurrentWrites(true),
		sftp.MaxConcurrentRequestsPerFile(64),
	)
	if err != nil {
		return err
	}
	defer client.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	dstFile, err := client.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	return su.uploadStrategy.Copy(srcFile, dstFile)
}
