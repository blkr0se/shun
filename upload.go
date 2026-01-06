package shun

import (
	"bytes"

	"github.com/pkg/sftp"
)

type Uploader interface {
	Upload(data []byte, path string) error
}

// Uploads data to a SSH server using SFTP.
func (r *shun) Upload(data []byte, path string) error {
	sftp, err := sftp.NewClient(r.client)
	if err != nil {
		return err
	}
	defer sftp.Close()

	srcFile := bytes.NewReader(data)
	dstFile, err := sftp.Create(path)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// write to file
	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		return err
	}
	return nil
}
