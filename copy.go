package shun

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/sftp"
	"github.com/schollz/progressbar/v3"
)

type RemoteFileCopier interface {
	Copy(dst io.Writer, src io.Reader) error
}

type AtOnce struct {
	written int64
}

// Tries to copy a file all at once.
func (ao *AtOnce) Copy(dst io.Writer, src io.Reader) error {
	written, err := io.Copy(dst, src)

	ao.written = written

	return err
}

type Chunked struct {
	Source      *os.File
	Destination *sftp.File
}

func (c *Chunked) defaultProgressBar(size int64) *progressbar.ProgressBar {
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

// Copies a file and displays progress to the tty.
func (c *Chunked) Copy(dst io.Writer, src io.Reader) error {
	// Create progress bar that will be displayed to the client.
	stats, err := c.Source.Stat()
	if err != nil {
		return err
	}
	bar := c.defaultProgressBar(stats.Size())

	// Upload a few chunks of 'src' at a time.
	// TODO(bia): Should use a totalWritten accumulator as the loop condition.
	for bar.State().CurrentBytes < float64(bar.State().Max) {
		written, err := io.CopyN(dst, src, 6*1000*1000)
		if err != nil && err != io.EOF {
			return err
		}

		// Increment progress bar that will be rendered on the next refresh.
		if err := bar.Add64(written); err != nil {
			return err
		}
	}

	return nil
}
