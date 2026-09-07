package main

import (
	"flag"
	"log/slog"
	"os"

	"goredis/internal/server"
)

func main() {
	port := flag.String("port", "6379", "TCP port to listen on")
	flag.Parse()

	s := server.New(":" + *port)
	if err := s.ListenAndServe(); err != nil {
		slog.Error("server exited", "err", err)
		os.Exit(1)
	}
}
