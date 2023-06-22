package main

import (
	"kings-corner/internal/grpc"
	"kings-corner/internal/storage"
)

func main() {
	s := storage.New()
	gServer := grpc.New(s)

	s.Run()
	gServer.Run()
}
