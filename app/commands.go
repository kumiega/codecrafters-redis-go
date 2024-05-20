package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Handler func([]Value, *Storage) Value

func (f Handler) Handle(args []Value, storage *Storage) Value {
	return f(args, storage)
}

var commands map[string]Handler

func init() {
	commands = map[string]Handler{
		"COMMAND": command,
		"HELP":    help,
		"INFO":    info,
		"ECHO":    echo,
		"PING":    ping,
		"SET":     set,
		"GET":     get,
	}
}

// Handles the incoming command.
func HandleCommand(value Value, storage *Storage) Value {
	command := strings.ToUpper(value.array[0].bulk)
	args := value.array[1:]

	handler, ok := commands[command]

	if !ok {
		fmt.Println("Invalid command", command)
		return Value{dataType: "string", str: ""}
	}

	return handler.Handle(args, storage)
}

// Support testing with redis-cli.
func command(args []Value, storage *Storage) Value {
	return Value{dataType: "string", str: "To see all available commands use 'help' command"}
}

// Provides a list of available commands.
func help(args []Value, storage *Storage) Value {
	commandList := make([]Value, 0, len(commands))

	for cmd := range commands {
		if cmd == "COMMAND" {
			continue
		}

		commandList = append(commandList, Value{dataType: "string", str: cmd})
	}

	return Value{dataType: "array", array: commandList}
}

func info(args []Value, storage *Storage) Value {
	info := ""
	info += "# Replication\n"
	info += "role:master\n"

	return Value{dataType: "bulk", bulk: info}
}

// Echoes back the input.
func echo(args []Value, storage *Storage) Value {
	if len(args) != 1 {
		return Value{dataType: "error", str: "ERR wrong number of arguments for 'echo' command"}
	}

	return Value{dataType: "bulk", bulk: args[0].bulk}
}

// Replies with PONG or echoes the provided argument.
func ping(args []Value, storage *Storage) Value {
	if len(args) == 0 {
		return Value{dataType: "string", str: "PONG"}
	}

	return Value{dataType: "string", str: args[0].bulk}
}

// Sets a value in the storage.
func set(args []Value, storage *Storage) Value {
	if len(args) != 2 && len(args) != 4 {
		return Value{dataType: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	if len(args) == 4 {
		if strings.ToUpper(args[2].bulk) != "PX" {
			return Value{dataType: "error", str: "ERR invalid argument for 'set' command"}
		}

		if _, err := strconv.Atoi(args[3].bulk); err != nil {
			return Value{dataType: "error", str: "ERR expiration should be integer"}
		}

		err := storage.SetWithExpiration(args[0].bulk, args[1].bulk, args[3].bulk)
		if err != nil {
			return Value{dataType: "error", str: "ERR could not set expiration time"}
		}

		return Value{dataType: "string", str: "OK"}
	}

	storage.Set(args[0].bulk, args[1].bulk)
	return Value{dataType: "string", str: "OK"}
}

// Retrieves a value from the storage.
func get(args []Value, storage *Storage) Value {
	if len(args) != 1 {
		return Value{dataType: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	value := storage.Get(args[0].bulk)

	if value == "" {
		return Value{dataType: "null"}
	}

	return Value{dataType: "bulk", bulk: value}
}
