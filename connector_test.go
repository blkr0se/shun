package shun_test

import (
	"net"
	"testing"
	"time"

	"github.com/blkr0se/shun"
	server "github.com/gliderlabs/ssh"
)

func startMockSSHServer(port string) error {
	server.Handle(func(s server.Session) {
		s.Write([]byte("pong"))
	})

	l, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		return err
	}

	go server.Serve(l, nil)
	time.Sleep(100 * time.Millisecond) // Give server time to start
	return nil
}

func TestConnect(t *testing.T) {
	// Start mock SSH server on port 2222
	if err := startMockSSHServer("2222"); err != nil {
		t.Fatalf("failed to start mock SSH server: %v", err)
	}

	sh := shun.NewShun("")
	if sh == nil {
		t.Fatal("expected shun instance, got nil")
	}

	params := shun.ConnectionParams{
		User: "me",
		Ip:   "127.0.0.1",
		Port: "2222",
	}

	if _, err := sh.Connect(params); err != nil {
		t.Fatalf("expected no error on connect, got %v", err)
	}
}
