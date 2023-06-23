package main

import (
	"kings-corner/internal/grpc"
	"kings-corner/internal/services"
	"kings-corner/internal/storage"
)

func main() {
	s := storage.New()

	boardService := services.NewBoardService(s.BoardRepository(), s.PlayerRepository())

	grpcGameService := grpc.NewGameService(boardService)
	grpcPlayerService := grpc.NewPlayerService(boardService)
	gServer := grpc.New(s, grpcGameService, grpcPlayerService)

	s.Run()
	gServer.Run()
}
