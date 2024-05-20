package main

import "flag"

type Context struct {
	Storage *Storage
	Replica *Replica
	Port    *int
}

type Replica struct {
	Role             string
	MasterReplOffset int
	MasterReplId     string
}

// Returns new context
func NewContext() *Context {
	port := flag.Int("port", 6379, "The port to listen on")
	replicaof := flag.String("replicaof", "", "Is the slave replica?")
	flag.Parse()

	storage := NewStore()

	role := "master"
	if *replicaof != "" {
		role = "slave"
	}

	replica := &Replica{
		Role:             role,
		MasterReplOffset: 0,
		MasterReplId:     GenerateRandomString(40),
	}

	ctx := &Context{
		Storage: storage,
		Replica: replica,
		Port:    port,
	}

	return ctx
}
