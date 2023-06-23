package grpc

import (
	"kings-corner/internal/deck"
	"kings-corner/internal/game"
	"kings-corner/internal/mappers"
	"kings-corner/pkg/pb"
)

type stream interface {
	Send(*pb.JoinResponse) error
}

func handleBoardListenEvent(
	board *game.Board,
	player game.Player,
	stream stream,
) error {
	err := stream.Send(&pb.JoinResponse{
		Board:        mappers.BoardDomainToPb(*board),
		PlayerId:     player.ID(),
		Hand:         buildPbHand(player.Hand()),
		IsPlayerTurn: isPlayerTurn(board, player.ID()),
	})
	if err != nil {
		return err
	}

	for {
		b := board.Listen()

		err = stream.Send(&pb.JoinResponse{
			Board:        mappers.BoardDomainToPb(b),
			PlayerId:     player.ID(),
			Hand:         buildPbHand(player.Hand()),
			IsPlayerTurn: isPlayerTurn(board, player.ID()),
		})
		if err != nil {
			break
		}
	}

	return err
}

func buildPbHand(hand []deck.Card) []*pb.Card {
	pbHand := []*pb.Card{}
	for _, c := range hand {
		pbHand = append(pbHand, mappers.CardDomainToPb(c))
	}

	return pbHand
}

func isPlayerTurn(board *game.Board, playerID string) bool {
	for i := range board.Players {
		if i == int(board.CurrentTurn) && playerID == board.Players[i].ID() {
			return true
		}
	}

	return false
}
