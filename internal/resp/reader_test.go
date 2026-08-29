// internal/resp/reader_test.go
package resp

import (
	"strings"
	"testing"
)

func TestReadArrayCommand(t *testing.T) {
	r := NewReader(strings.NewReader("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"))
	v, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"SET", "foo", "bar"}
	if len(v.Array) != len(want) {
		t.Fatalf("got %d items, want %d", len(v.Array), len(want))
	}
	for i, w := range want {
		if v.Array[i].Str != w {
			t.Fatalf("item %d: got %q, want %q", i, v.Array[i].Str, w)
		}
	}
}

func TestReadArrayCommandNullBulk(t *testing.T) {
	r := NewReader(strings.NewReader("*1\r\n$-1\r\n"))
	v, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v.Array) != 1 || !v.Array[0].Null {
		t.Fatalf("got %+v, want one null bulk string", v.Array)
	}
}

func TestReadInlineCommand(t *testing.T) {
	r := NewReader(strings.NewReader("PING\r\n"))
	v, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v.Array) != 1 || v.Array[0].Str != "PING" {
		t.Fatalf("got %+v, want [PING]", v.Array)
	}
}

func TestReadInlineCommandMultipleArgs(t *testing.T) {
	r := NewReader(strings.NewReader("SET foo bar\r\n"))
	v, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"SET", "foo", "bar"}
	for i, w := range want {
		if v.Array[i].Str != w {
			t.Fatalf("item %d: got %q, want %q", i, v.Array[i].Str, w)
		}
	}
}

func TestReadMultipleCommandsPipelined(t *testing.T) {
	r := NewReader(strings.NewReader("PING\r\n*1\r\n$4\r\nPING\r\n"))
	for i := 0; i < 2; i++ {
		v, err := r.Read()
		if err != nil {
			t.Fatalf("read %d: unexpected error: %v", i, err)
		}
		if v.Array[0].Str != "PING" {
			t.Fatalf("read %d: got %q, want PING", i, v.Array[0].Str)
		}
	}
}
