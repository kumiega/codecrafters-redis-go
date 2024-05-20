package main

import (
	"fmt"
	"io"
	"net"
)

// Handles incomming connection
func HandleConnection(conn net.Conn, storage *Storage) {
	defer conn.Close()

	fmt.Println("New connection established")

	for {
		resp := NewRespReader(conn)
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

		response := HandleCommand(value, storage)

		writer := NewRespWriter(conn)
		writer.Write(response)
	}
}
