package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	_STRING  = '+'
	_ERROR   = '-'
	_INTEGER = ':'
	_BULK    = '$'
	_ARRAY   = '*'
)

type Value struct {
	dataType string
	str      string
	num      int
	bulk     string
	array    []Value
}

type Resp struct {
	reader *bufio.Reader
}

// Returns RESP Reader
func NewRespReader(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// Reads bytes and determines content
func (r *Resp) Read() (Value, error) {
	dataType, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch dataType {
	case _ARRAY:
		return r.readArray()
	case _BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(dataType))
		return Value{}, nil
	}
}

// Reads each line
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

// Reads bytes and converts to integer
func (r *Resp) readInteger() (x int, n int, err error) {
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

// Reads bytes as RESP array
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.dataType = "array"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	v.array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		v.array = append(v.array, val)
	}

	return v, nil
}

// Reads bytes as RESP bulk
func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.dataType = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	r.readLine()

	return v, nil
}
