package commands

import (
	"strings"

	"goredis/internal/resp"
	"goredis/internal/storage"
)

type Handler func(store *storage.Store, args []string) resp.Value

var table = map[string]Handler{
	"PING":   ping,
	"ECHO":   echo,
	"GET":    get,
	"SET":    set,
	"SETNX":  setnx,
	"DEL":    del,
	"EXISTS": exists,
	"INCR":   incr,
	"DECR":   decr,
	"APPEND": appendCmd,
	"STRLEN": strlen,
	"MGET":   mget,
	"MSET":   mset,
}

// Dispatch runs the command named by args[0] against store and returns the
// RESP reply. args must be non-empty.
func Dispatch(store *storage.Store, args []string) resp.Value {
	name := strings.ToUpper(args[0])
	handler, ok := table[name]
	if !ok {
		return resp.ErrorValue("ERR unknown command '" + args[0] + "'")
	}
	return handler(store, args[1:])
}

func ping(store *storage.Store, args []string) resp.Value {
	if len(args) == 0 {
		return resp.SimpleStringValue("PONG")
	}
	return resp.BulkStringValue(args[0])
}

func echo(store *storage.Store, args []string) resp.Value {
	if len(args) != 1 {
		return resp.ErrorValue("ERR wrong number of arguments for 'echo' command")
	}
	return resp.BulkStringValue(args[0])
}
