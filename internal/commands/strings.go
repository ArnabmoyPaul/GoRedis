package commands

import (
	"goredis/internal/resp"
	"goredis/internal/storage"
)

func get(store *storage.Store, args []string) resp.Value {
	if len(args) != 1 {
		return resp.ErrorValue("ERR wrong number of arguments for 'get' command")
	}
	v, ok := store.Get(args[0])
	if !ok {
		return resp.NullBulkString()
	}
	return resp.BulkStringValue(v)
}

func set(store *storage.Store, args []string) resp.Value {
	if len(args) != 2 {
		return resp.ErrorValue("ERR wrong number of arguments for 'set' command")
	}
	store.Set(args[0], args[1])
	return resp.SimpleStringValue("OK")
}

func setnx(store *storage.Store, args []string) resp.Value {
	if len(args) != 2 {
		return resp.ErrorValue("ERR wrong number of arguments for 'setnx' command")
	}
	if store.SetNX(args[0], args[1]) {
		return resp.IntegerValue(1)
	}
	return resp.IntegerValue(0)
}

func del(store *storage.Store, args []string) resp.Value {
	if len(args) < 1 {
		return resp.ErrorValue("ERR wrong number of arguments for 'del' command")
	}
	return resp.IntegerValue(int64(store.Del(args...)))
}

func exists(store *storage.Store, args []string) resp.Value {
	if len(args) < 1 {
		return resp.ErrorValue("ERR wrong number of arguments for 'exists' command")
	}
	return resp.IntegerValue(int64(store.Exists(args...)))
}

func incr(store *storage.Store, args []string) resp.Value {
	if len(args) != 1 {
		return resp.ErrorValue("ERR wrong number of arguments for 'incr' command")
	}
	n, err := store.Incr(args[0])
	if err != nil {
		return resp.ErrorValue("ERR " + err.Error())
	}
	return resp.IntegerValue(n)
}

func decr(store *storage.Store, args []string) resp.Value {
	if len(args) != 1 {
		return resp.ErrorValue("ERR wrong number of arguments for 'decr' command")
	}
	n, err := store.Decr(args[0])
	if err != nil {
		return resp.ErrorValue("ERR " + err.Error())
	}
	return resp.IntegerValue(n)
}

func appendCmd(store *storage.Store, args []string) resp.Value {
	if len(args) != 2 {
		return resp.ErrorValue("ERR wrong number of arguments for 'append' command")
	}
	return resp.IntegerValue(int64(store.Append(args[0], args[1])))
}

func strlen(store *storage.Store, args []string) resp.Value {
	if len(args) != 1 {
		return resp.ErrorValue("ERR wrong number of arguments for 'strlen' command")
	}
	return resp.IntegerValue(int64(store.StrLen(args[0])))
}

func mget(store *storage.Store, args []string) resp.Value {
	if len(args) < 1 {
		return resp.ErrorValue("ERR wrong number of arguments for 'mget' command")
	}
	values := store.MGet(args...)
	items := make([]resp.Value, len(values))
	for i, v := range values {
		if v == nil {
			items[i] = resp.NullBulkString()
		} else {
			items[i] = resp.BulkStringValue(*v)
		}
	}
	return resp.ArrayValue(items...)
}

func mset(store *storage.Store, args []string) resp.Value {
	if len(args) < 2 || len(args)%2 != 0 {
		return resp.ErrorValue("ERR wrong number of arguments for 'mset' command")
	}
	pairs := make(map[string]string, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		pairs[args[i]] = args[i+1]
	}
	store.MSet(pairs)
	return resp.SimpleStringValue("OK")
}
