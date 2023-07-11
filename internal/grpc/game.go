package grpc

import (
	"context"
	"log"

	"kings-corner/internal/game"
	"kings-corner/internal/pubsub"
	"kings-corner/internal/services"
	"kings-corner/pkg/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GameService struct {
	boardService *services.BoardService
	boardPubsub  pubsub.PubSub[game.Board]

	pb.UnimplementedGameServiceServer
}

func NewGameService(b *services.BoardService, boardPubsub pubsub.PubSub[game.Board]) *GameService {
	return &GameService{b, boardPubsub, pb.UnimplementedGameServiceServer{}}
}

func (s *GameService) Create(_ *pb.CreateRequest, stream pb.GameService_CreateServer) error {
	board, err := s.boardService.Create()
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	log.Printf("board %s created", board.ID)

	event := NewEventHandler(board.Players[0], stream)
	err = event.Handle(board, s.boardPubsub.Subscribe(board.Channel()))
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}

func (s *GameService) Begin(_ context.Context, req *pb.BeginRequest) (*pb.BeginResponse, error) {
	err := s.boardService.Run(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	log.Printf("board %s began", req.Id)

	return &pb.BeginResponse{}, nil
}
