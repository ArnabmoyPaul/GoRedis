package server

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func startTestServer(t *testing.T) net.Addr {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	s := New(ln.Addr().String())
	go s.Serve(ln)
	t.Cleanup(func() { ln.Close() })
	return ln.Addr()
}

func TestServerPing(t *testing.T) {
	addr := startTestServer(t)
	conn, err := net.DialTimeout("tcp", addr.String(), time.Second)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("*1\r\n$4\r\nPING\r\n")); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	reply, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if reply != "+PONG\r\n" {
		t.Fatalf("got %q, want %q", reply, "+PONG\r\n")
	}
}

func TestServerSetGet(t *testing.T) {
	addr := startTestServer(t)
	conn, err := net.DialTimeout("tcp", addr.String(), time.Second)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	if _, err := conn.Write([]byte("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	line, _ := r.ReadString('\n')
	if line != "+OK\r\n" {
		t.Fatalf("got %q, want +OK\\r\\n", line)
	}

	if _, err := conn.Write([]byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n")); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	header, _ := r.ReadString('\n')
	if header != "$3\r\n" {
		t.Fatalf("got %q, want $3\\r\\n", header)
	}
	body, _ := r.ReadString('\n')
	if body != "bar\r\n" {
		t.Fatalf("got %q, want bar\\r\\n", body)
	}
}

func TestServerInlinePing(t *testing.T) {
	addr := startTestServer(t)
	conn, err := net.DialTimeout("tcp", addr.String(), time.Second)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("PING\r\n")); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	reply, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if reply != "+PONG\r\n" {
		t.Fatalf("got %q, want +PONG\\r\\n", reply)
	}
}

func TestServerPipelinedCommands(t *testing.T) {
	addr := startTestServer(t)
	conn, err := net.DialTimeout("tcp", addr.String(), time.Second)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	// Two PINGs sent back-to-back with no read in between.
	if _, err := conn.Write([]byte("*1\r\n$4\r\nPING\r\n*1\r\n$4\r\nPING\r\n")); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	r := bufio.NewReader(conn)
	for i := 0; i < 2; i++ {
		reply, err := r.ReadString('\n')
		if err != nil {
			t.Fatalf("read %d failed: %v", i, err)
		}
		if reply != "+PONG\r\n" {
			t.Fatalf("read %d: got %q, want +PONG\\r\\n", i, reply)
		}
	}
}
