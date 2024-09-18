package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	jobpb "github.com/rudeigerc/go-grpc-protovalidate/examples/gen/go/job/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func init() {
	flag.Parse()
}

func main() {
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := jobpb.NewJobServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	job := jobpb.Job{DisplayName: "!!!"}
	r, err := c.CreateJob(ctx, &jobpb.CreateJobRequest{Parent: "namespaces/foo", Job: &job})
	if err != nil {
		s := status.Convert(err)
		log.Printf("Error: %v", err.Error())
		for _, d := range s.Details() {
			switch info := d.(type) {
			case *errdetails.ErrorInfo:
				log.Printf("Error info: %s", info)
			case *errdetails.BadRequest:
				log.Printf("Bad request: %s", info)
			default:
				log.Printf("Unexpected type: %s", info)
			}
		}
		os.Exit(1)
	}
	log.Printf("Job: %v", r)
}
