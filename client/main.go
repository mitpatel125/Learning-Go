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

	resp, err := client.Echo(ctx, &pb.EchoRequest{
		Message: "grpc is now clicking",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Client received:", resp.Message)

	listResp, err := client.GetMessages(ctx, &pb.GetMessagesRequest{Limit: 5})
	if err != nil{
		log.Fatal(err)
	}


	log.Println("Latest messages:")
	for _, m := range listResp.Messages{
		t := time.Unix(m.CreatedUnix, 0).Format(time.RFC3339)
		log.Printf("- #%d [%s] %s", m.Id, t, m.Content)
	}

	if len(listResp.Messages) > 0{
		id:= listResp.Messages[0].Id


		readResp, err := client.MarkMessageAsRead(ctx, &pb.MarkMessageAsReadRequest{
			MessageId: int64(id),
		})
		if err != nil{
			log.Fatal(err)
		}
		log.Printf(
			"marked as read : #%d (%v)",
			readResp.Messages.Id,
			readResp.Messages.Read,
		)
		getResp, err :=  client.GetMessageByID(ctx, &pb.GetMessageByIDRequest{
			MessageId: int64(id),
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(
			"fetched by id: #%d | %s | read=%v",
    	getResp.Message.Id,
    	getResp.Message.Content,
    	getResp.Message.Read,)
	


	}

}
