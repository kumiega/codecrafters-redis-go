package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

type Context struct {
	Storage *Storage
	Replica *string
	Port    *int
}

func main() {
	storage := NewStore()

	port := flag.Int("port", 6379, "The port to listen on")
	replica := flag.String("replicaof", "", "Is the slave replica?")
	flag.Parse()

	addr := fmt.Sprintf("0.0.0.0:%d", *port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to bind to port %d: %v\n", *port, err)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Printf("Server starting on port %d\n", *port)

	ctx := &Context{
		Storage: storage,
		Replica: replica,
		Port:    port,
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go HandleConnection(conn, ctx)
	}
}
