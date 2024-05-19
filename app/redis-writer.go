package main

import (
	"io"
	"strconv"
)

type RedisRespWriter struct {
	writer io.Writer
}

func NewRedisRespWriter(w io.Writer) *RedisRespWriter {
	return &RedisRespWriter{writer: w}
}

func (w *RedisRespWriter) Write(v RedisValue) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (v RedisValue) Marshal() []byte {
	switch v.dataType {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshallNull()
	case "error":
		return v.marshallError()
	default:
		return []byte{}
	}
}

func (v RedisValue) marshalArray() []byte {
	len := len(v.array)
	var bytes []byte
	bytes = append(bytes, REDIS_ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes
}

func (v RedisValue) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, REDIS_BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v RedisValue) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, REDIS_STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v RedisValue) marshallNull() []byte {
	return []byte("$-1\r\n")
}

func (v RedisValue) marshallError() []byte {
	var bytes []byte
	bytes = append(bytes, REDIS_ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}
