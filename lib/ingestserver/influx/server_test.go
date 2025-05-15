package influx

import (
	"net"
	"testing"
	"time"
)

func TestServerConnTimeout(t *testing.T) {
	addr := "localhost:18189"
	timeout := 2 * time.Second
	server := &Server{}
	go func() {
		if err := Start(addr, timeout); err != nil {
			t.Fatalf("failed to start server: %v", err)
		}
	}()
	defer server.Stop()

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	time.Sleep(timeout + 500*time.Millisecond)

	_, err = conn.Write([]byte("test_metric value=42\n"))
	if err == nil {
		t.Fatalf("expected connection to be closed after timeout")
	}
}

func TestServerNoTimeout(t *testing.T) {
	addr := "localhost:18190"
	server := &Server{}
	go func() {
		if err := Start(addr, 0); err != nil {
			t.Fatalf("failed to start server: %v", err)
		}
	}()
	defer server.Stop()

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	time.Sleep(2 * time.Second)

	_, err = conn.Write([]byte("test_metric value=42\n"))
	if err != nil {
		t.Fatalf("expected connection to remain open: %v", err)
	}
}

func (s *Server) Stop() {
	if s.someChannel != nil { // Replace someChannel with the actual channel name
		close(s.someChannel)
	}
}
