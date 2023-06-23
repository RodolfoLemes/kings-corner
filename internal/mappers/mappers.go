package mappers

import (
	"kings-corner/internal/deck"
	"kings-corner/internal/game"
	"kings-corner/pkg/pb"
)

func BoardDomainToPb(b game.Board) *pb.Board {
	playerIDs := []string{}
	for _, p := range b.Players {
		playerIDs = append(playerIDs, p.ID())
	}

	pbFields := []*pb.Board_Field{}
	for i, fieldCards := range b.Field {
		pbCards := []*pb.Card{}
		for _, c := range fieldCards {
			pbCards = append(pbCards, CardDomainToPb(c))
		}
		pbFields = append(pbFields, &pb.Board_Field{
			FieldNumber: uint32(i),
			Cards:       pbCards,
		})
	}

	return &pb.Board{
		Id:          "example",
		PlayerIds:   playerIDs,
		Fields:      pbFields,
		CurrentTurn: uint32(b.CurrentTurn),
	}
}

func CardDomainToPb(c deck.Card) *pb.Card {
	return &pb.Card{
		Suit: uint32(c.Suit),
		Rank: uint32(c.Rank),
	}
}
