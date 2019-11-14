package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"time"

	pb "github.com/pgermishuys/esgrpc/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address     = "localhost:2113"
	defaultName = "Event Store GRPC"
)

func main() {
	// Set up a connection to the server.
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStreamsClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	log.Printf(name)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	readReq := &pb.ReadReq{
		Options: &pb.ReadReq_Options{
			CountOptions: &pb.ReadReq_Options_Count{
				Count: 4096,
			},
			FilterOptions: &pb.ReadReq_Options_NoFilter{
				NoFilter: &pb.ReadReq_Empty{},
			},
			ReadDirection: pb.ReadReq_Options_Forwards,
			ResolveLinks:  true,
			StreamOptions: &pb.ReadReq_Options_Stream{
				Stream: &pb.ReadReq_Options_StreamOptions{
					StreamName: "foo",
					RevisionOptions: &pb.ReadReq_Options_StreamOptions_Start{
						Start: &pb.ReadReq_Empty{},
					},
				},
			},
		},
	}

	result, err := c.Read(ctx, readReq)
	if err != nil {
		log.Fatalf("could not read: %v", err)
	}
	log.Printf("%+v", result)
}
