package grpc

import (
	"log"
	"net"

	"kings-corner/internal/storage"
	"kings-corner/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server interface {
	Run()
}

type gRPCServer struct {
	storage storage.Storage
}

func New(s storage.Storage) Server {
	return &gRPCServer{s}
}

func (*gRPCServer) Run() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Could not connected: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGameServiceServer(grpcServer, &GameService{})
	pb.RegisterPlayerServiceServer(grpcServer, &PlayerService{})

	reflection.Register(grpcServer)

	log.Println("gRPC Server Running...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Could not server: %v", err)
	}
}
