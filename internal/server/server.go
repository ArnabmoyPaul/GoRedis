package server

import (
	"errors"
	"io"
	"log/slog"
	"net"

	"goredis/internal/commands"
	"goredis/internal/resp"
	"goredis/internal/storage"
)

type Server struct {
	addr  string
	store *storage.Store
}

func New(addr string) *Server {
	return &Server{addr: addr, store: storage.New()}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	return s.Serve(ln)
}

// Serve accepts connections on ln until Accept fails (e.g. the listener is
// closed). Exposed separately from ListenAndServe so tests can bind to an
// ephemeral port and discover the real address before serving.
func (s *Server) Serve(ln net.Listener) error {
	defer ln.Close()
	slog.Info("goredis listening", "addr", ln.Addr().String())
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	reader := resp.NewReader(conn)
	for {
		cmd, err := reader.Read()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				slog.Warn("connection read error", "err", err)
			}
			return
		}
		if len(cmd.Array) == 0 {
			continue
		}
		args := make([]string, len(cmd.Array))
		for i, v := range cmd.Array {
			args[i] = v.Str
		}
		reply := commands.Dispatch(s.store, args)
		if _, err := conn.Write(reply.Encode()); err != nil {
			return
		}
	}
}
