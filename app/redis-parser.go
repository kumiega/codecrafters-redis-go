package main

import (
	"errors"
	"fmt"
	"strings"
)

type RedisCommand struct {
	Command string
	Args    []string
}

func ProcessRedisClientCommand(bytes []byte) ([]byte, error) {
	command, err := getRedisCommand(bytes)

	if err != nil {
		return []byte{}, err
	}

	resp := handleRedisResponse(command)

	return resp, nil
}

func getRedisCommand(bytes []byte) (RedisCommand, error) {
	var (
		command string
		args    []string
	)

	parts := strings.Split(string(bytes), "\r\n")

	if len(parts) < 2 || !strings.Contains(parts[0], "*") {
		return RedisCommand{}, errors.New("invalid command content")
	}

	for idx, arg := range parts[2:] {
		if arg == "" {
			continue
		}

		if idx == 0 {
			command = strings.ToUpper(arg)
			continue
		}

		if strings.HasPrefix(arg, "$") {
			continue
		}

		args = append(args, arg)
	}

	return RedisCommand{
		Command: command,
		Args:    args,
	}, nil
}

func handleRedisResponse(command RedisCommand) []byte {
	var response string

	for _, arg := range command.Args {
		fmt.Println(arg)
	}

	switch command.Command {
	case "PING":
		response = RedisSimpleString("PONG")
	case "ECHO":
		response = RedisBulkString(command.Args[0])
	}

	return []byte(response)
}
