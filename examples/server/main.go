package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	go_grpc_protovalidate "github.com/rudeigerc/go-grpc-protovalidate"
	jobpb "github.com/rudeigerc/go-grpc-protovalidate/examples/gen/go/job/v1"
	"google.golang.org/grpc"
)

type server struct {
	jobpb.UnimplementedJobServiceServer
}

func (s *server) CreateJob(ctx context.Context, req *jobpb.CreateJobRequest) (*jobpb.Job, error) {
	return &jobpb.Job{
		Name:        uuid.Must(uuid.NewRandom()).String(),
		DisplayName: req.GetJob().GetDisplayName(),
	}, nil
}

var (
	port = flag.Int("port", 50051, "The server port")
)

func init() {
	flag.Parse()
}

func main() {
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("failed to create validator: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			go_grpc_protovalidate.UnaryServerInterceptor(go_grpc_protovalidate.WithValidator(validator)),
		),
		grpc.ChainStreamInterceptor(
			go_grpc_protovalidate.StreamServerInterceptor(go_grpc_protovalidate.WithValidator(validator)),
		),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	jobpb.RegisterJobServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
