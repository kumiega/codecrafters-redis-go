package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	ctx := NewContext()

	addr := fmt.Sprintf("0.0.0.0:%d", *ctx.Port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to bind to port %d: %v\n", *ctx.Port, err)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Printf("Server starting on port %d\n", *ctx.Port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go HandleConnection(conn, ctx)
	}
}
