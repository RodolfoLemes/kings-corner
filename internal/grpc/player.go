package grpc

import (
	"context"

	"kings-corner/internal/deck"
	"kings-corner/internal/game"
	"kings-corner/internal/services"
	pb "kings-corner/pkg/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return status.Errorf(codes.Internal, err.Error())
	}

	err = handleBoardListenEvent(board, *player, stream)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
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
	case pb.PlayRequest_MOVE:
		if len(req.MoveTurn.FieldCardLevel) != 2 {
			return nil, status.Error(
				codes.InvalidArgument,
				"invalid field card level, should have 2 elements",
			)
		}

		turn = &game.MoveTurn{
			FieldLevel: [2]uint8{
				uint8(req.MoveTurn.FieldCardLevel[0]),
				uint8(req.MoveTurn.FieldCardLevel[1]),
			},
			MoveToFieldLevel: uint8(req.MoveTurn.MoveToFieldLevel),
		}
	case pb.PlayRequest_PASS:
		turn = &game.PassTurn{}
	}

	err := s.boardService.Play(req.GameId, req.PlayerId, turn)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.PlayResponse{}, nil
}
