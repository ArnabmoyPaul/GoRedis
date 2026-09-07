package storage

import "testing"

func TestSetGet(t *testing.T) {
	s := New()
	s.Set("foo", "bar")
	v, ok := s.Get("foo")
	if !ok || v != "bar" {
		t.Fatalf("got (%q, %v), want (bar, true)", v, ok)
	}
}

func TestGetMissing(t *testing.T) {
	s := New()
	if _, ok := s.Get("missing"); ok {
		t.Fatalf("expected ok=false for missing key")
	}
}

func TestSetNX(t *testing.T) {
	s := New()
	if !s.SetNX("foo", "bar") {
		t.Fatalf("expected SetNX to succeed on new key")
	}
	if s.SetNX("foo", "baz") {
		t.Fatalf("expected SetNX to fail on existing key")
	}
	v, _ := s.Get("foo")
	if v != "bar" {
		t.Fatalf("got %q, want bar (unchanged)", v)
	}
}

func TestDel(t *testing.T) {
	s := New()
	s.Set("foo", "bar")
	s.Set("baz", "qux")
	n := s.Del("foo", "missing", "baz")
	if n != 2 {
		t.Fatalf("got %d deleted, want 2", n)
	}
	if _, ok := s.Get("foo"); ok {
		t.Fatalf("expected foo to be deleted")
	}
}

func TestExists(t *testing.T) {
	s := New()
	s.Set("foo", "bar")
	if n := s.Exists("foo", "missing", "foo"); n != 2 {
		t.Fatalf("got %d, want 2", n)
	}
}

func TestIncrOnMissingKey(t *testing.T) {
	s := New()
	n, err := s.Incr("counter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Fatalf("got %d, want 1", n)
	}
}

func TestIncrOnExistingKey(t *testing.T) {
	s := New()
	s.Set("counter", "10")
	n, err := s.Incr("counter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 11 {
		t.Fatalf("got %d, want 11", n)
	}
}

func TestIncrOnNonIntegerValue(t *testing.T) {
	s := New()
	s.Set("foo", "bar")
	if _, err := s.Incr("foo"); err == nil {
		t.Fatalf("expected error for non-integer value")
	}
}

func TestDecr(t *testing.T) {
	s := New()
	s.Set("counter", "10")
	n, err := s.Decr("counter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 9 {
		t.Fatalf("got %d, want 9", n)
	}
}

func TestAppend(t *testing.T) {
	s := New()
	if n := s.Append("greeting", "Hello "); n != 6 {
		t.Fatalf("got %d, want 6", n)
	}
	if n := s.Append("greeting", "World"); n != 11 {
		t.Fatalf("got %d, want 11", n)
	}
	v, _ := s.Get("greeting")
	if v != "Hello World" {
		t.Fatalf("got %q, want %q", v, "Hello World")
	}
}

func TestStrLen(t *testing.T) {
	s := New()
	s.Set("foo", "hello")
	if n := s.StrLen("foo"); n != 5 {
		t.Fatalf("got %d, want 5", n)
	}
	if n := s.StrLen("missing"); n != 0 {
		t.Fatalf("got %d, want 0", n)
	}
}

func TestMGet(t *testing.T) {
	s := New()
	s.Set("foo", "bar")
	out := s.MGet("foo", "missing")
	if out[0] == nil || *out[0] != "bar" {
		t.Fatalf("got %v, want bar", out[0])
	}
	if out[1] != nil {
		t.Fatalf("got %v, want nil", out[1])
	}
}

func TestMSet(t *testing.T) {
	s := New()
	s.MSet(map[string]string{"a": "1", "b": "2"})
	va, _ := s.Get("a")
	vb, _ := s.Get("b")
	if va != "1" || vb != "2" {
		t.Fatalf("got (%q, %q), want (1, 2)", va, vb)
	}
}
