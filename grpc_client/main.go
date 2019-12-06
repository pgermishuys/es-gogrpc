package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"time"

	"github.com/gofrs/uuid"
	pb "github.com/pgermishuys/es-gogrpc/protos"
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
	client := pb.NewStreamsClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	log.Printf(name)
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	// readRequest := &pb.ReadReq{
	// 	Options: &pb.ReadReq_Options{
	// 		CountOptions: &pb.ReadReq_Options_Count{
	// 			Count: 4096,
	// 		},
	// 		FilterOptions: &pb.ReadReq_Options_NoFilter{
	// 			NoFilter: &pb.ReadReq_Empty{},
	// 		},
	// 		ReadDirection: pb.ReadReq_Options_Forwards,
	// 		ResolveLinks:  true,
	// 		StreamOptions: &pb.ReadReq_Options_Stream{
	// 			Stream: &pb.ReadReq_Options_StreamOptions{
	// 				StreamName: "foo",
	// 				RevisionOptions: &pb.ReadReq_Options_StreamOptions_Start{
	// 					Start: &pb.ReadReq_Empty{},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	id, _ := uuid.NewV4()
	idAsBinary, _ := id.MarshalBinary()
	appendRequestHeader := &pb.AppendReq{
		Content: &pb.AppendReq_Options_{
			Options: &pb.AppendReq_Options{
				StreamName: "foo",
				Id:         idAsBinary,
				ExpectedStreamRevision: &pb.AppendReq_Options_Any{
					Any: &pb.AppendReq_Empty{},
				},
			},
		},
	}
	appendRequest := &pb.AppendReq{
		Content: &pb.AppendReq_ProposedMessage_{
			ProposedMessage: &pb.AppendReq_ProposedMessage{
				Id:   idAsBinary,
				Data: []byte{},
			},
		},
	}
	log.Printf("%+v", id)
	log.Printf("%+v", idAsBinary)

	// result, err := c.Read(ctx, readRequest)
	// log.Printf("%+v", result)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	streamClient, err := client.Append(ctx)
	if err != nil {
		log.Fatalf("could not get streams client: %v", err)
	}
	//send header
	if err := streamClient.Send(appendRequestHeader); err != nil {
		log.Fatalf("Can not send header %+v", err)
	}
	// result, err := streamClient.CloseAndRecv()
	// if err != nil {
	// 	log.Fatalf("Can not send header %+v", err)
	// }
	// uuidFromBytes, err := uuid.FromBytes(result.Id)
	// log.Printf("%+v", uuidFromBytes.String())
	//send request
	if err = streamClient.Send(appendRequest); err != nil {
		log.Printf("sending request error: %v", err)
	}
	// log.Printf("%+v", result)
}
