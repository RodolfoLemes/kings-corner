package grpc

import (
	"context"

	"kings-corner/internal/deck"
	"kings-corner/internal/game"
	"kings-corner/internal/services"
	pb "kings-corner/pkg/pb"
)

type PlayerService struct {
	boardService *services.BoardService
	pb.UnimplementedPlayerServiceServer
}

func NewPlayerService(b *services.BoardService) *PlayerService {
	return &PlayerService{b, pb.UnimplementedPlayerServiceServer{}}
}

func (s *PlayerService) Join(req *pb.JoinRequest, stream pb.PlayerService_JoinServer) error {
	board, player, err := s.boardService.Join(req.GameId)
	if err != nil {
		return err
	}

	return handleBoardListenEvent(board, *player, stream)
}

func (s *PlayerService) Play(_ context.Context, req *pb.PlayRequest) (*pb.PlayResponse, error) {
	var turn game.Turn
	switch req.Turn {
	case pb.PlayRequest_CARD:
		turn = &game.CardTurn{
			FieldLevel: uint8(req.CardTurn.FieldLevel),
			Card: deck.Card{
				Suit: deck.Suit(req.CardTurn.Card.Suit),
				Rank: deck.Rank(req.CardTurn.Card.Rank),
			},
		}
	}

	return nil, s.boardService.Play(req.GameId, req.PlayerId, turn)
}
