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
