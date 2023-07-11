package grpc

import (
	"context"
	"log"

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

	log.Printf("player %s joined in board %s", player.ID(), board.ID)

	err = handleBoardListenEvent(board, player, stream)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}

func (s *PlayerService) Play(_ context.Context, req *pb.PlayRequest) (*pb.PlayResponse, error) {
	var turn game.Turn
	switch req.TurnMode {
	case pb.PlayRequest_CARD:
		pbTurn := req.Turn.(*pb.PlayRequest_CardTurn_).CardTurn

		turn = &game.CardTurn{
			FieldLevel: uint8(pbTurn.FieldLevel),
			Card: deck.Card{
				Suit: deck.Suit(pbTurn.Card.Suit),
				Rank: deck.Rank(pbTurn.Card.Rank),
			},
		}
	case pb.PlayRequest_MOVE:
		pbTurn := req.Turn.(*pb.PlayRequest_MoveTurn_).MoveTurn
		if len(pbTurn.FieldCardLevel) != 2 {
			return nil, status.Error(
				codes.InvalidArgument,
				"invalid field card level, should have 2 elements",
			)
		}

		turn = &game.MoveTurn{
			FieldLevel: [2]uint8{
				uint8(pbTurn.FieldCardLevel[0]),
				uint8(pbTurn.FieldCardLevel[1]),
			},
			MoveToFieldLevel: uint8(pbTurn.MoveToFieldLevel),
		}
	case pb.PlayRequest_PASS:
		turn = &game.PassTurn{}
	}

	err := s.boardService.Play(req.GameId, req.PlayerId, turn)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	log.Printf("player %s played", req.PlayerId)

	return &pb.PlayResponse{}, nil
}
