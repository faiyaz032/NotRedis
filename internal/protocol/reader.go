package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type RESPReader struct {
	reader *bufio.Reader
}

func NewRESPReader(rd io.Reader) *RESPReader {
	return &RESPReader{
		reader: bufio.NewReader(rd),
	}
}

func (r *RESPReader) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n++
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *RESPReader) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *RESPReader) ReadValue() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch string(_type) {
	case Array:
		return r.readArray()
	case BulkString:
		return r.readBulk()
	case SimpleString:
		return r.readSimpleString()
	case Integer:
		return r.readIntegerType()
	case Error:
		return r.readError()
	default:
		fmt.Printf("Unknown type: %v\n", string(_type))
		return Value{}, fmt.Errorf("unknown type: %v", string(_type))
	}
}

func (r *RESPReader) readArray() (Value, error) {
	v := Value{Type: Array}

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	if len == -1 {
		v.Null = true
		return v, nil
	}

	v.Array = make([]Value, len)
	for i := 0; i < len; i++ {
		val, err := r.ReadValue()
		if err != nil {
			return v, err
		}
		v.Array[i] = val
	}

	return v, nil
}

func (r *RESPReader) readBulk() (Value, error) {
	v := Value{Type: BulkString}

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	if len == -1 {
		v.Null = true
		return v, nil
	}

	bulk := make([]byte, len)
	r.reader.Read(bulk)
	v.Bulk = string(bulk)

	// read trailing crlf
	r.readLine()

	return v, nil
}

func (r *RESPReader) readSimpleString() (Value, error) {
	v := Value{Type: SimpleString}
	line, _, err := r.readLine()
	if err != nil {
		return v, err
	}
	v.Str = string(line)
	return v, nil
}

func (r *RESPReader) readIntegerType() (Value, error) {
	v := Value{Type: Integer}
	x, _, err := r.readInteger()
	if err != nil {
		return v, err
	}
	v.Num = x
	return v, nil
}

func (r *RESPReader) readError() (Value, error) {
	v := Value{Type: Error}
	line, _, err := r.readLine()
	if err != nil {
		return v, err
	}
	v.Str = string(line)
	return v, nil
}
