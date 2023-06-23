package grpc

import (
	"context"

	"kings-corner/internal/services"
	"kings-corner/pkg/pb"
)

type GameService struct {
	boardService *services.BoardService

	pb.UnimplementedGameServiceServer
}

func NewGameService(b *services.BoardService) *GameService {
	return &GameService{b, pb.UnimplementedGameServiceServer{}}
}

func (s *GameService) Create(_ *pb.CreateRequest, stream pb.GameService_CreateServer) error {
	board, err := s.boardService.Create()
	if err != nil {
		return err
	}

	return handleBoardListenEvent(board, board.Players[0], stream)
}

func (s *GameService) Begin(_ context.Context, req *pb.BeginRequest) (*pb.BeginResponse, error) {
	return nil, s.boardService.Run(req.Id)
}
