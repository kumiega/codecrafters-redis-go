package main

import "strconv"

func RedisSimpleString(str string) string {
	return "+" + str + "\r\n"
}

func RedisSimpleError(str string) string {
	return "-" + str + "\r\n"
}

func RedisSignedInteger(integer int, negative bool) string {
	var sign string

	if negative {
		sign = "-"
	} else {
		sign = "+"
	}

	return ":" + sign + strconv.Itoa(integer) + "\r\n"
}

func RedisInteger(integer int, negative bool) string {
	return ":" + strconv.Itoa(integer) + "\r\n"
}

func RedisBulkString(str string) string {
	if str == "" {
		return "$0\r\n\r\n"
	}

	length := strconv.Itoa(len(str))

	return "$" + length + "\r\n" + str + "\r\n"
}
