package main

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)

type RedisAof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewRedisAof(path string) (*RedisAof, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	aof := &RedisAof{
		file: file,
		rd:   bufio.NewReader(file),
	}

	go func() {
		aof.mu.Lock()
		aof.file.Sync()
		aof.mu.Unlock()
		time.Sleep(time.Second)
	}()

	return aof, nil
}

func (aof *RedisAof) Read(fn func(value RedisValue)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	aof.file.Seek(0, io.SeekStart)

	reader := NewRedisRespReader(aof.file)

	for {
		value, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		fn(value)
	}

	return nil
}

func (aof *RedisAof) Write(value RedisValue) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (aof *RedisAof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}
