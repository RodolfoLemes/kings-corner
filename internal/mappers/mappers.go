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

	pbPlayers := []*pb.Board_Player{}
	for _, p := range b.Players {
		pbPlayers = append(pbPlayers, &pb.Board_Player{
			Id:   p.ID(),
			Hand: uint32(len(p.Hand())),
		})
	}

	return &pb.Board{
		Id:          b.ID,
		Players:     pbPlayers,
		Fields:      pbFields,
		CurrentTurn: uint32(b.CurrentTurn),
		IsStarted:   b.IsStarted,
	}
}

func CardDomainToPb(c deck.Card) *pb.Card {
	return &pb.Card{
		Suit: uint32(c.Suit),
		Rank: uint32(c.Rank),
	}
}

func CardPbToDomain(pb *pb.Card) deck.Card {
	return deck.Card{
		Suit: deck.Suit(pb.Suit),
		Rank: deck.Rank(pb.Rank),
	}
}
