package cli

import (
	"context"
	"log"
	"os"

	pb "kings-corner/pkg/pb"

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

	joinAction := &JoinAction{playerClient}
	createAction := &CreateAction{gameClient, playerClient}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "player",
				Subcommands: []*cli.Command{
					{
						Name:   "join",
						Action: joinAction.Action,
					},
				},
			},
			{
				Name: "game",
				Subcommands: []*cli.Command{
					{
						Name:   "create",
						Action: createAction.Action,
					},
					{
						Name: "begin",
						Action: func(cCtx *cli.Context) error {
							gameID := cCtx.Args().First()

							_, err := gameClient.Begin(context.TODO(), &pb.BeginRequest{
								Id: gameID,
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
