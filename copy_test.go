package shun_test

import (
	"math"
	"os"
	"strings"
	"testing"

	"github.com/blkr0se/shun"
)

var (
	fakeSrc = "shun_copy_test_fake_src"
	fakeDst = strings.ReplaceAll(fakeSrc, "src", "dst")
	fourGb  = 4 * math.Pow(2, 30)
)

func createFakeFile(name string) (*os.File, error) {
	data := make([]byte, int(fourGb))
	f, err := os.CreateTemp("", name)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(f.Name(), data, os.ModeAppend); err != nil {
		return nil, err
	}

	return f, nil
}

func TestChunkedCopy(t *testing.T) {
	// Bootstrap
	src, err := createFakeFile(fakeSrc)
	if err != nil {
		t.Fatalf("error creating fake src: %s", err)
	}

	dst, err := createFakeFile(fakeDst)
	if err != nil {
		t.Fatalf("error creating fake dst: %s", err)
	}

	c := &shun.Chunked{}
	if err := c.Copy(src, dst); err != nil {
		t.Fatalf("error copying src to dst: %s", err)
	}

	// Cleanup
	if err := os.Remove(src.Name()); err != nil {
		t.Fatalf("error removing src: %s", err)
	}

	if err := os.Remove(dst.Name()); err != nil {
		t.Fatalf("error removing dst: %s", err)
	}
}
