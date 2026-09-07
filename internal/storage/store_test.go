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
