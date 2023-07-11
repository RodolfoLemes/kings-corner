package grpc

import (
	"log"

	"kings-corner/internal/game"
	"kings-corner/internal/mappers"
	"kings-corner/pkg/pb"
)

type stream interface {
	Send(*pb.JoinResponse) error
}

type EventHandler struct {
	player game.Player
	stream stream
}

func NewEventHandler(player game.Player, stream stream) *EventHandler {
	return &EventHandler{player, stream}
}

func (e *EventHandler) Handle(board *game.Board, listener <-chan game.Board) error {
	err := e.sendToStream(board)
	if err != nil {
		return err
	}

	for {
		b := <-listener

		err = e.sendToStream(&b)
		if err != nil {
			log.Println(err)
			break
		}

	}

	return err
}

func (e EventHandler) isPlayerTurn(board *game.Board) bool {
	for i := range board.Players {
		if i == int(board.CurrentTurn) && e.player.ID() == board.Players[i].ID() {
			return true
		}
	}

	return false
}

func (e EventHandler) sendToStream(board *game.Board) error {
	message := &pb.JoinResponse{
		Board:        mappers.BoardDomainToPb(*board),
		PlayerId:     e.player.ID(),
		Hand:         mappers.BuildPbHand(e.player.Hand()),
		IsPlayerTurn: e.isPlayerTurn(board),
	}

	err := e.stream.Send(message)

	log.Printf("Sending message to %s: %+v \n", board.ID, message)

	return err
}
