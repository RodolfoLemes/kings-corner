package grpc

import (
	"context"
	"fmt"

	"kings-corner/internal/deck"
	"kings-corner/internal/game"
	pb "kings-corner/pkg/pb"
)

type PlayerService struct {
	pb.UnimplementedPlayerServiceServer
}

func (*PlayerService) Join(_ *pb.JoinRequest, stream pb.PlayerService_JoinServer) error {
	p := game.NewPlayer()
	board.Join(p)

	isPlayerTurn := false
	playerIDs := []string{}
	for _, p := range board.Players {
		playerIDs = append(playerIDs, p.ID())
	}

	pbFields := []*pb.Board_Field{}
	for i, fieldCards := range board.Field {
		pbCards := []*pb.Card{}
		for _, c := range fieldCards {
			pbCards = append(pbCards, &pb.Card{
				Suit: uint32(c.Suit),
				Rank: uint32(c.Rank),
			})
		}
		pbFields = append(pbFields, &pb.Board_Field{
			FieldNumber: uint32(i),
			Cards:       pbCards,
		})
	}

	pbHand := []*pb.Card{}
	for _, c := range p.Hand() {
		pbHand = append(pbHand, &pb.Card{
			Suit: uint32(c.Suit),
			Rank: uint32(c.Rank),
		})
	}

	err := stream.Send(&pb.JoinResponse{
		Board: &pb.Board{
			Id:          "example",
			PlayerIds:   playerIDs,
			Fields:      pbFields,
			CurrentTurn: uint32(board.CurrentTurn),
		},
		PlayerId:     p.ID(),
		Hand:         pbHand,
		IsPlayerTurn: isPlayerTurn,
	})

	for {
		b := board.Listen()

		isPlayerTurn = false
		playerIDs := []string{}
		for i, p := range b.Players {
			playerIDs = append(playerIDs, p.ID())
			if i == int(b.CurrentTurn) {
				isPlayerTurn = true
			}
		}

		pbFields := []*pb.Board_Field{}
		for i, fieldCards := range b.Field {
			pbCards := []*pb.Card{}
			for _, c := range fieldCards {
				pbCards = append(pbCards, &pb.Card{
					Suit: uint32(c.Suit),
					Rank: uint32(c.Rank),
				})
			}
			pbFields = append(pbFields, &pb.Board_Field{
				FieldNumber: uint32(i),
				Cards:       pbCards,
			})
		}

		pbHand := []*pb.Card{}
		for _, c := range p.Hand() {
			pbHand = append(pbHand, &pb.Card{
				Suit: uint32(c.Suit),
				Rank: uint32(c.Rank),
			})
		}

		err = stream.Send(&pb.JoinResponse{
			Board: &pb.Board{
				Id:          "example",
				PlayerIds:   playerIDs,
				Fields:      pbFields,
				CurrentTurn: uint32(b.CurrentTurn),
			},
			PlayerId:     p.ID(),
			Hand:         pbHand,
			IsPlayerTurn: isPlayerTurn,
		})
		if err != nil {
			break
		}
	}

	return err
}

func (*PlayerService) Play(_ context.Context, req *pb.PlayRequest) (*pb.PlayResponse, error) {
	var player game.Player
	for i := range board.Players {
		if board.Players[i].ID() == req.PlayerId {
			player = board.Players[i]
			break
		}
	}

	fmt.Printf("req.PlayerId: %v\n", req.PlayerId)
	fmt.Printf("board.Players: %v\n", board.Players)
	fmt.Printf("player: %v\n", player)

	switch req.Turn {
	case pb.PlayRequest_CARD:
		ct := &game.CardTurn{
			FieldLevel: uint8(req.CardTurn.FieldLevel),
			Card: deck.Card{
				Suit: deck.Suit(req.CardTurn.Card.Suit),
				Rank: deck.Rank(req.CardTurn.Card.Rank),
			},
		}

		player.PlayTurn(ct)
	}

	return nil, nil
}
