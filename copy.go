package shun

import (
	"fmt"
	"io"
	"os"
	"time"

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
	written int64

	tracker *progressbar.ProgressBar
}

func (c *Chunked) defaultProgressBar(size int64) {
	c.tracker = progressbar.NewOptions64(
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
	stats, err := src.(*os.File).Stat()
	if err != nil {
		return err
	}

	c.defaultProgressBar(stats.Size())
	c.beginCopy(dst, src)

	return nil
}

// Begin copying a few chunks of 'src' at a time.
// TODO(bia): Should use a totalWritten accumulator as the loop condition.
func (c *Chunked) beginCopy(dst io.Writer, src io.Reader) error {
	for c.tracker.State().CurrentBytes < float64(c.tracker.State().Max) {
		if err := c.copyChunk(dst, src); err != nil {
			return err
		}

		// Increment progress bar that will be rendered on the next refresh.
		if err := c.incrementProgress(); err != nil {
			return err
		}
	}

	return nil
}

// Copies a file and displays progress to the tty.
func (c *Chunked) copyChunk(dst io.Writer, src io.Reader) error {
	written, err := io.CopyN(dst, src, 6*1000*1000)
	if err != nil && err != io.EOF {
		return err
	}

	c.written = written
	return nil
}

func (c *Chunked) incrementProgress() error {
	return c.tracker.Add64(c.written)
}
