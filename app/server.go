package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	storage := NewStore()

	l, err := net.Listen("tcp", "0.0.0.0:6379")

	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(" server starting on port 6379")

	defer l.Close()
	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err.Error())
			continue
		}

		go HandleConnection(conn, storage)
	}
}
