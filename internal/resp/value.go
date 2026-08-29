// internal/resp/value.go
package resp

import (
	"fmt"
	"strconv"
)

type Type byte

const (
	SimpleString Type = '+'
	Error        Type = '-'
	Integer      Type = ':'
	BulkString   Type = '$'
	Array        Type = '*'
)

type Value struct {
	Type  Type
	Str   string
	Num   int64
	Array []Value
	Null  bool
}

func SimpleStringValue(s string) Value { return Value{Type: SimpleString, Str: s} }
func ErrorValue(s string) Value        { return Value{Type: Error, Str: s} }
func IntegerValue(n int64) Value       { return Value{Type: Integer, Num: n} }
func BulkStringValue(s string) Value   { return Value{Type: BulkString, Str: s} }
func NullBulkString() Value            { return Value{Type: BulkString, Null: true} }
func ArrayValue(items ...Value) Value  { return Value{Type: Array, Array: items} }
func NullArray() Value                 { return Value{Type: Array, Null: true} }

func (v Value) Encode() []byte {
	switch v.Type {
	case SimpleString:
		return []byte("+" + v.Str + "\r\n")
	case Error:
		return []byte("-" + v.Str + "\r\n")
	case Integer:
		return []byte(":" + strconv.FormatInt(v.Num, 10) + "\r\n")
	case BulkString:
		if v.Null {
			return []byte("$-1\r\n")
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v.Str), v.Str))
	case Array:
		if v.Null {
			return []byte("*-1\r\n")
		}
		buf := []byte(fmt.Sprintf("*%d\r\n", len(v.Array)))
		for _, item := range v.Array {
			buf = append(buf, item.Encode()...)
		}
		return buf
	default:
		return nil
	}
}
