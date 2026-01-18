package shun

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/sftp"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/crypto/ssh"
)

// Uploader is responsible for uploading content to a remote private server.
type Uploader interface {

	// Upload writes bytes to a file in dst on a remote server.
	Upload(src string, dst string) error
}

type SftpUploader struct {
	sshClient *ssh.Client

	showProgress bool
}

func NewSftpUploader(s *ssh.Client) *SftpUploader {
	return &SftpUploader{sshClient: s}
}

func (su *SftpUploader) UseProgressBar(v bool) error {
	su.showProgress = v
	return nil
}

// Upload writes bytes to a file on a remote server using SFTP.
func (su *SftpUploader) Upload(src string, dst string) error {
	sftp, err := sftp.NewClient(
		su.sshClient,
		sftp.UseConcurrentReads(true),
		sftp.UseConcurrentWrites(true),
		sftp.MaxConcurrentRequestsPerFile(64),
	)
	if err != nil {
		return err
	}
	defer sftp.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	dstFile, err := sftp.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if su.showProgress {
		stats, _ := srcFile.Stat()
		bar := su.tempProgressBar(stats.Size())
		for bar.State().CurrentBytes < float64(bar.State().Max) {
			written, err := io.CopyN(dstFile, srcFile, 6*1000*1000)
			if err != nil {
				if err != io.EOF {
					return err
				}
			}

			if err := bar.Add64(written); err != nil {
				return err
			}
		}
	} else {
		_, err := io.Copy(dstFile, srcFile)
		return err
	}

	return nil
}

func (su *SftpUploader) tempProgressBar(size int64) *progressbar.ProgressBar {
	return progressbar.NewOptions64(
		size,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowTotalBytes(true),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetWidth(50),
	)
}
