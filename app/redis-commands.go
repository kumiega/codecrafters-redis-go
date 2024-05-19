package main

import "sync"

var RedisCommands = map[string]func([]RedisValue) RedisValue{
	"COMMAND": command,
	"ECHO":    echo,
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func command(args []RedisValue) RedisValue {
	return RedisValue{dataType: "string", str: "Look elsewhere else for Redis commands."}
}

func echo(args []RedisValue) RedisValue {
	if len(args) != 1 {
		return RedisValue{dataType: "error", str: "ERR wrong number of arguments for 'echo' command"}
	}

	return RedisValue{dataType: "bulk", bulk: args[0].bulk}
}

func ping(args []RedisValue) RedisValue {
	if len(args) == 0 {
		return RedisValue{dataType: "string", str: "PONG"}
	}

	return RedisValue{dataType: "string", str: args[0].bulk}
}

var SET_STORAGE = map[string]string{}
var SET_STORAGE_MUTEX = sync.RWMutex{}

func set(args []RedisValue) RedisValue {
	if len(args) != 2 {
		return RedisValue{dataType: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SET_STORAGE_MUTEX.Lock()
	SET_STORAGE[key] = value
	SET_STORAGE_MUTEX.Unlock()

	return RedisValue{dataType: "string", str: "OK"}
}

func get(args []RedisValue) RedisValue {
	if len(args) != 1 {
		return RedisValue{dataType: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	SET_STORAGE_MUTEX.RLock()
	value, ok := SET_STORAGE[key]
	SET_STORAGE_MUTEX.RUnlock()

	if !ok {
		return RedisValue{dataType: "null"}
	}

	return RedisValue{dataType: "bulk", bulk: value}
}

var HSET_STORAGE = map[string]map[string]string{}
var HSET_STORAGE_MUTEX = sync.RWMutex{}

func hset(args []RedisValue) RedisValue {
	if len(args) != 3 {
		return RedisValue{dataType: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSET_STORAGE_MUTEX.Lock()
	if _, ok := HSET_STORAGE[hash]; !ok {
		HSET_STORAGE[hash] = map[string]string{}
	}
	HSET_STORAGE[hash][key] = value
	HSET_STORAGE_MUTEX.Unlock()

	return RedisValue{dataType: "string", str: "OK"}
}

func hget(args []RedisValue) RedisValue {
	if len(args) != 2 {
		return RedisValue{dataType: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSET_STORAGE_MUTEX.RLock()
	value, ok := HSET_STORAGE[hash][key]
	HSET_STORAGE_MUTEX.RUnlock()

	if !ok {
		return RedisValue{dataType: "null"}
	}

	return RedisValue{dataType: "bulk", bulk: value}
}

func hgetall(args []RedisValue) RedisValue {
	if len(args) != 1 {
		return RedisValue{dataType: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	hash := args[0].bulk

	HSET_STORAGE_MUTEX.RLock()
	value, ok := HSET_STORAGE[hash]
	HSET_STORAGE_MUTEX.RUnlock()

	if !ok {
		return RedisValue{dataType: "null"}
	}

	values := []RedisValue{}

	for k, v := range value {
		values = append(values, RedisValue{dataType: "bulk", bulk: k})
		values = append(values, RedisValue{dataType: "bulk", bulk: v})
	}

	return RedisValue{dataType: "array", array: values}
}
