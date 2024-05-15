package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("Redis server starting on port 6379")

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("New connection established")

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from connection:", err)
			}
			return
		}

		fmt.Printf("Recived: \n%q\n", string(buf[:n]))

		response, err := ProcessRedisClientCommand(buf[:n])
		if err != nil {
			fmt.Println("Error with command processing:", err)
			return
		}

		fmt.Printf("Responded: \n%q\n", string(response))

		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
	}
}
