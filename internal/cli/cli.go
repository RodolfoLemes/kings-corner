package cli

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	pb "kings-corner/pkg/pb"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func Run() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer conn.Close()

	playerClient := pb.NewPlayerServiceClient(conn)
	gameClient := pb.NewGameServiceClient(conn)

	// TODO - refactor this mess
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "player",
				Subcommands: []*cli.Command{
					{
						Name: "join",
						Action: func(cCtx *cli.Context) error {
							stream, err := playerClient.Join(context.TODO(), &pb.JoinRequest{
								GameId: "example",
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

								fmt.Printf("player.stream.Board: %v\n", stream.Board)
								fmt.Printf("player.stream.PlayerId: %v\n", stream.PlayerId)
								fmt.Printf("player.stream.Hand: %v\n", stream.Hand)
								fmt.Printf("player.stream.IsPlayerTurn: %v\n", stream.IsPlayerTurn)

								if !stream.IsPlayerTurn {
									continue
								}

								prompt := promptui.Select{
									Label: "Choose action",
									Items: []string{"Play a Card", "Move", "Skip"},
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
										selections = append(selections, c.String())
									}

									prompt := promptui.Select{
										Label: "Choose card to play",
										Items: selections,
									}

									_, result, err := prompt.Run()
									if err != nil {
										fmt.Printf("Prompt failed %v\n", err)
										continue
									}

									playedCard := &pb.Card{}
									for _, c := range stream.Hand {
										if c.String() == result {
											playedCard = c
										}
									}

									prompt = promptui.Select{
										Label: "Choose a field to play",
										Items: []uint8{1, 2, 3, 4, 5, 6, 7, 8},
									}

									_, result, err = prompt.Run()
									if err != nil {
										fmt.Printf("Prompt failed %v\n", err)
										continue
									}

									fieldLevel, err := strconv.Atoi(result)
									if err != nil {
										fmt.Printf("Prompt failed %v\n", err)
										continue
									}

									resp, err := playerClient.Play(cCtx.Context, &pb.PlayRequest{
										Turn: pb.PlayRequest_CARD,
										CardTurn: &pb.PlayRequest_CardTurn{
											FieldLevel: uint32(fieldLevel - 1),
											Card:       playedCard,
										},
										PlayerId: stream.PlayerId,
									})
									if err != nil {
										fmt.Println(err)
										continue
									}

									fmt.Printf("resp: %v\n", resp)
								}
							}

							return nil
						},
					},
				},
			},
			{
				Name: "game",
				Subcommands: []*cli.Command{
					{
						Name: "begin",
						Action: func(cCtx *cli.Context) error {
							_, err := gameClient.Begin(context.TODO(), &pb.BeginRequest{
								Id: "example",
							})
							if err != nil {
								return err
							}

							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
