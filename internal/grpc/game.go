package grpc

import (
	"kings-corner/internal/deck"
	"kings-corner/internal/game"
	"kings-corner/pkg/pb"
)

var board game.Board

func init() {
	d := deck.New()

	board = game.New(d)
}

type GameService struct {
	pb.UnimplementedGameServiceServer
}

func (*GameService) Create(_ *pb.CreateRequest, stream pb.GameService_CreateServer) error {
	return nil
}

func (*GameService) Begin(req *pb.BeginRequest, stream pb.GameService_BeginServer) error {
	// id doesn't matter right now
	board.Run()

	var err error
	for {
		b := board.Listen()
		playerIDs := []string{}
		for _, p := range b.Players {
			playerIDs = append(playerIDs, p.ID())
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

		err := stream.Send(&pb.BeginResponse{
			Board: &pb.Board{
				Id:          "example",
				PlayerIds:   playerIDs,
				Fields:      pbFields,
				CurrentTurn: uint32(b.CurrentTurn),
			},
		})
		if err != nil {
			break
		}
	}

	return err
}
