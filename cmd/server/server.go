package main

import (
	"kings-corner/internal/game"
	"kings-corner/internal/grpc"
	"kings-corner/internal/pubsub"
	"kings-corner/internal/services"
	"kings-corner/internal/storage"
)

func main() {
	s := storage.New()

	boardPs := pubsub.New[game.Board]()

	boardService := services.NewBoardService(s.BoardRepository(), s.PlayerRepository(), boardPs)

	grpcGameService := grpc.NewGameService(boardService, boardPs)
	grpcPlayerService := grpc.NewPlayerService(boardService, boardPs)
	gServer := grpc.New(s, grpcGameService, grpcPlayerService)

	s.Run()
	gServer.Run()
}
