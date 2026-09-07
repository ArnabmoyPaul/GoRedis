package commands

import (
	"testing"

	"goredis/internal/resp"
	"goredis/internal/storage"
)

func TestPing(t *testing.T) {
	s := storage.New()
	v := Dispatch(s, []string{"PING"})
	if v.Type != resp.SimpleString || v.Str != "PONG" {
		t.Fatalf("got %+v, want PONG", v)
	}
}

func TestPingWithMessage(t *testing.T) {
	s := storage.New()
	v := Dispatch(s, []string{"PING", "hello"})
	if v.Type != resp.BulkString || v.Str != "hello" {
		t.Fatalf("got %+v, want hello", v)
	}
}

func TestEcho(t *testing.T) {
	s := storage.New()
	v := Dispatch(s, []string{"ECHO", "hello"})
	if v.Type != resp.BulkString || v.Str != "hello" {
		t.Fatalf("got %+v, want hello", v)
	}
}

func TestUnknownCommand(t *testing.T) {
	s := storage.New()
	v := Dispatch(s, []string{"NOPE"})
	if v.Type != resp.Error {
		t.Fatalf("got %+v, want error", v)
	}
}

func TestCommandNameCaseInsensitive(t *testing.T) {
	s := storage.New()
	v := Dispatch(s, []string{"ping"})
	if v.Type != resp.SimpleString || v.Str != "PONG" {
		t.Fatalf("got %+v, want PONG", v)
	}
}

func TestSetAndGet(t *testing.T) {
	s := storage.New()
	Dispatch(s, []string{"SET", "foo", "bar"})
	v := Dispatch(s, []string{"GET", "foo"})
	if v.Type != resp.BulkString || v.Str != "bar" {
		t.Fatalf("got %+v, want bar", v)
	}
}

func TestGetMissingReturnsNull(t *testing.T) {
	s := storage.New()
	v := Dispatch(s, []string{"GET", "missing"})
	if !v.Null {
		t.Fatalf("got %+v, want null bulk string", v)
	}
}

func TestSetNXCommand(t *testing.T) {
	s := storage.New()
	v := Dispatch(s, []string{"SETNX", "foo", "bar"})
	if v.Num != 1 {
		t.Fatalf("got %+v, want 1", v)
	}
	v = Dispatch(s, []string{"SETNX", "foo", "baz"})
	if v.Num != 0 {
		t.Fatalf("got %+v, want 0", v)
	}
}

func TestDelCommand(t *testing.T) {
	s := storage.New()
	Dispatch(s, []string{"SET", "foo", "bar"})
	v := Dispatch(s, []string{"DEL", "foo", "missing"})
	if v.Num != 1 {
		t.Fatalf("got %+v, want 1", v)
	}
}

func TestExistsCommand(t *testing.T) {
	s := storage.New()
	Dispatch(s, []string{"SET", "foo", "bar"})
	v := Dispatch(s, []string{"EXISTS", "foo", "missing"})
	if v.Num != 1 {
		t.Fatalf("got %+v, want 1", v)
	}
}

func TestIncrDecrCommands(t *testing.T) {
	s := storage.New()
	v := Dispatch(s, []string{"INCR", "counter"})
	if v.Num != 1 {
		t.Fatalf("got %+v, want 1", v)
	}
	v = Dispatch(s, []string{"DECR", "counter"})
	if v.Num != 0 {
		t.Fatalf("got %+v, want 0", v)
	}
}

func TestAppendCommand(t *testing.T) {
	s := storage.New()
	v := Dispatch(s, []string{"APPEND", "greeting", "Hello"})
	if v.Num != 5 {
		t.Fatalf("got %+v, want 5", v)
	}
}

func TestStrLenCommand(t *testing.T) {
	s := storage.New()
	Dispatch(s, []string{"SET", "foo", "hello"})
	v := Dispatch(s, []string{"STRLEN", "foo"})
	if v.Num != 5 {
		t.Fatalf("got %+v, want 5", v)
	}
}

func TestMsetAndMget(t *testing.T) {
	s := storage.New()
	Dispatch(s, []string{"MSET", "a", "1", "b", "2"})
	v := Dispatch(s, []string{"MGET", "a", "b", "missing"})
	if len(v.Array) != 3 {
		t.Fatalf("got %d items, want 3", len(v.Array))
	}
	if v.Array[0].Str != "1" || v.Array[1].Str != "2" || !v.Array[2].Null {
		t.Fatalf("got %+v", v.Array)
	}
}

func TestMsetOddArgsIsError(t *testing.T) {
	s := storage.New()
	v := Dispatch(s, []string{"MSET", "a", "1", "b"})
	if v.Type != resp.Error {
		t.Fatalf("got %+v, want error", v)
	}
}
