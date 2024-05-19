package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	REDIS_STRING  = '+'
	REDIS_ERROR   = '-'
	REDIS_INTEGER = ':'
	REDIS_BULK    = '$'
	REDIS_ARRAY   = '*'
)

type RedisValue struct {
	dataType string
	str      string
	num      int
	bulk     string
	array    []RedisValue
}

type RedisResp struct {
	reader *bufio.Reader
}

func NewRedisRespReader(rd io.Reader) *RedisResp {
	return &RedisResp{reader: bufio.NewReader(rd)}
}

func (r *RedisResp) Read() (RedisValue, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return RedisValue{}, err
	}

	switch _type {
	case REDIS_ARRAY:
		return r.readArray()
	case REDIS_BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return RedisValue{}, nil
	}
}

func (r *RedisResp) readLine() (line []byte, n int, err error) {
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

func (r *RedisResp) readInteger() (x int, n int, err error) {
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

func (r *RedisResp) readArray() (RedisValue, error) {
	v := RedisValue{}
	v.dataType = "array"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	v.array = make([]RedisValue, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		v.array = append(v.array, val)
	}

	return v, nil
}

func (r *RedisResp) readBulk() (RedisValue, error) {
	v := RedisValue{}

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
