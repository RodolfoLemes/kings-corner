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

type JoinAction struct {
	playerClient pb.PlayerServiceClient
}

func (j *JoinAction) Action(cCtx *cli.Context) error {
	gameID := cCtx.Args().First()

	stream, err := j.playerClient.Join(cCtx.Context, &pb.JoinRequest{
		GameId: gameID,
	})
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

		fmt.Printf("PlayerID => %s", stream.PlayerId)

		if !stream.Board.IsStarted {
			continue
		}

		j.printBoard(stream.Board)
		j.printHand(stream.Hand)

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

			_, err = j.playerClient.Play(cCtx.Context, &pb.PlayRequest{
				TurnMode: pb.PlayRequest_CARD,
				Turn: &pb.PlayRequest_CardTurn_{
					CardTurn: &pb.PlayRequest_CardTurn{
						FieldLevel: uint32(fieldLevel),
						Card:       playedCard,
					},
				},
				PlayerId: stream.PlayerId,
				GameId:   gameID,
			})
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}

	return err
}

func (j *JoinAction) printBoard(board *pb.Board) {
	fmt.Println()
	for i, f := range board.Fields {
		cards := []deck.Card{}
		for _, c := range f.Cards {
			cards = append(cards, mappers.CardPbToDomain(c))
		}

		printCard := ""
		for _, c := range cards {
			printCard += fmt.Sprintf("%s-", c.String())
		}

		fmt.Printf("%d -> [%s] \n", i+1, printCard)
	}
}

func (j *JoinAction) printHand(pbCards []*pb.Card) {
	fmt.Println()

	printCard := ""
	for _, c := range pbCards {
		card := mappers.CardPbToDomain(c)

		printCard += fmt.Sprintf("%s-", card.String())
	}
	fmt.Printf("HAND -> [%s] \n", printCard)
}
