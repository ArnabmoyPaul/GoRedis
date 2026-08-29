// internal/resp/value_test.go
package resp

import "testing"

func TestEncodeSimpleString(t *testing.T) {
	got := string(SimpleStringValue("OK").Encode())
	want := "+OK\r\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestEncodeError(t *testing.T) {
	got := string(ErrorValue("ERR bad").Encode())
	want := "-ERR bad\r\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestEncodeInteger(t *testing.T) {
	got := string(IntegerValue(42).Encode())
	want := ":42\r\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestEncodeBulkString(t *testing.T) {
	got := string(BulkStringValue("foo").Encode())
	want := "$3\r\nfoo\r\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestEncodeNullBulkString(t *testing.T) {
	got := string(NullBulkString().Encode())
	want := "$-1\r\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestEncodeArray(t *testing.T) {
	got := string(ArrayValue(BulkStringValue("SET"), BulkStringValue("foo"), BulkStringValue("bar")).Encode())
	want := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestEncodeNullArray(t *testing.T) {
	got := string(NullArray().Encode())
	want := "*-1\r\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
