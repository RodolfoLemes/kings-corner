package cli

import (
	"fmt"
	"io"
	"log"

	"kings-corner/internal/deck"
	"kings-corner/internal/mappers"
	"kings-corner/pkg/pb"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

type CreateAction struct {
	gameClient   pb.GameServiceClient
	playerClient pb.PlayerServiceClient
}

func (ca *CreateAction) Action(cCtx *cli.Context) error {
	stream, err := ca.gameClient.Create(cCtx.Context, &pb.CreateRequest{})
	if err != nil {
		return err
	}

	for {
		stream, err := stream.Recv()
		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}

		if err == io.EOF {
			break
		}

		fmt.Printf("GameID => %s", stream.Board.Id)
		fmt.Printf("PlayerID => %s", stream.PlayerId)
		if !stream.Board.IsStarted {
			continue
		}

		ca.printBoard(stream.Board)
		ca.printHand(stream.Hand)

		if !stream.IsPlayerTurn {
			continue
		}

		prompt := promptui.Select{
			Label: "Choose action",
			Items: []string{"Play a Card", "Move a Card", "Skip"},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			continue
		}

		switch result {
		case "Play a Card":
			selections := []string{}
			for _, c := range stream.Hand {
				selections = append(selections, mappers.CardPbToDomain(c).String())
			}

			prompt := promptui.Select{
				Label: "Choose card to play",
				Items: selections,
			}

			index, _, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				continue
			}

			playedCard := stream.Hand[index]

			prompt = promptui.Select{
				Label: "Choose a field to play",
				Items: []uint8{1, 2, 3, 4, 5, 6, 7, 8},
			}

			index, _, err = prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				continue
			}

			fieldLevel := index

			resp, err := ca.playerClient.Play(cCtx.Context, &pb.PlayRequest{
				Turn: pb.PlayRequest_CARD,
				CardTurn: &pb.PlayRequest_CardTurn{
					FieldLevel: uint32(fieldLevel),
					Card:       playedCard,
				},
				PlayerId: stream.PlayerId,
				GameId:   stream.Board.Id,
			})
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Printf("resp: %v\n", resp)
		}
	}

	return err
}

func (CreateAction) printBoard(board *pb.Board) {
	fmt.Println()
	for i, f := range board.Fields {
		cards := []deck.Card{}
		for _, c := range f.Cards {
			cards = append(cards, mappers.CardPbToDomain(c))
		}

		printCard := ""
		for _, c := range cards {
			printCard += fmt.Sprintf("%s |", c.String())
		}

		fmt.Printf("%d -> [%s] \n", i+1, printCard)
	}
}

func (CreateAction) printHand(pbCards []*pb.Card) {
	printCard := ""
	for _, c := range pbCards {
		card := mappers.CardPbToDomain(c)

		printCard += fmt.Sprintf("%s |", card.String())
	}
	fmt.Printf("HAND -> [%s] \n", printCard)
}
