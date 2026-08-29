// internal/resp/reader.go
package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Reader struct {
	br *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{br: bufio.NewReader(r)}
}

// Read reads one command from the stream and always returns it as an Array
// of BulkStrings, whether the client sent a RESP array or an inline command.
func (r *Reader) Read() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	if len(line) == 0 {
		return r.Read()
	}
	if line[0] != '*' {
		return ArrayValue(inlineToValues(line)...), nil
	}
	count, err := strconv.Atoi(line[1:])
	if err != nil {
		return Value{}, fmt.Errorf("resp: invalid array header %q: %w", line, err)
	}
	if count < 0 {
		return NullArray(), nil
	}
	items := make([]Value, 0, count)
	for i := 0; i < count; i++ {
		item, err := r.readBulkString()
		if err != nil {
			return Value{}, err
		}
		items = append(items, item)
	}
	return ArrayValue(items...), nil
}

func (r *Reader) readBulkString() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	if len(line) == 0 || line[0] != '$' {
		return Value{}, fmt.Errorf("resp: expected bulk string, got %q", line)
	}
	n, err := strconv.Atoi(line[1:])
	if err != nil {
		return Value{}, fmt.Errorf("resp: invalid bulk string header %q: %w", line, err)
	}
	if n < 0 {
		return NullBulkString(), nil
	}
	buf := make([]byte, n+2) // payload + trailing CRLF
	if _, err := io.ReadFull(r.br, buf); err != nil {
		return Value{}, err
	}
	return BulkStringValue(string(buf[:n])), nil
}

func (r *Reader) readLine() (string, error) {
	line, err := r.br.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func inlineToValues(line string) []Value {
	fields := strings.Fields(line)
	values := make([]Value, len(fields))
	for i, f := range fields {
		values[i] = BulkStringValue(f)
	}
	return values
}
