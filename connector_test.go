package shun_test

import (
	"net"
	"testing"
	"time"

	"github.com/blkr0se/shun"
	server "github.com/gliderlabs/ssh"
)

func startMockSSHServer(port string) error {
	l, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		return err
	}

	go server.Serve(l, nil)
	time.Sleep(100 * time.Millisecond) // Give server time to start
	return nil
}

func TestConnect(t *testing.T) {
	port := "2222"
	// Start mock SSH server on port 2222
	if err := startMockSSHServer(port); err != nil {
		t.Fatalf("failed to start mock SSH server: %v", err)
	}

	params := &shun.SshConnector{
		User: "me",
		Ip:   "127.0.0.1",
		Port: "2222",
	}

	if _, err := params.Connect(); err != nil {
		t.Fatalf("expected no error on connect, got %v", err)
	}
}
