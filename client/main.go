package main

import (
	"context"
	"log"
	"time"

	pb "grpc-echo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewEchoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.UpperCase(ctx, &pb.EchoRequest{
		Message: "grpc is now clicking",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Client received:", resp.Message)
}
