package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	aof, err := NewRedisAof("redis.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Redis server starting on port 6379")

	defer aof.Close()

	aof.Read(func(value RedisValue) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := RedisCommands[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	defer l.Close()
	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn, aof)
	}
}

func handleConnection(conn net.Conn, aof *RedisAof) {
	defer conn.Close()

	fmt.Println("New connection established")

	for {
		resp := NewRedisRespReader(conn)
		value, err := resp.Read()

		if err == io.EOF {
			fmt.Println("Client disconnected")
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		if value.dataType != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewRedisRespWriter(conn)

		handler, ok := RedisCommands[command]

		if !ok {
			fmt.Println("Invalid command", command)
			writer.Write(RedisValue{dataType: "string", str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
