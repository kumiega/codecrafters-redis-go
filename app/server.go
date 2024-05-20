package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	storage := NewStore()

	port := flag.Int("port", 6379, "Port to listen on")
	flag.Parse()

	addr := fmt.Sprintf("0.0.0.0:%d", *port)

	l, err := net.Listen("tcp", addr)

	if err != nil {
		fmt.Printf("Failed to bind to port %d\n", *port)
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Server starting on port %d\n", *port)

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
