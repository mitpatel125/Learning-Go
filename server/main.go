package main

import (
	"context"
	"log"
	"net"
	"strings"

	pb "grpc-echo/proto"
	"google.golang.org/grpc"
)

type echoServer struct{
	pb.UnimplementedEchoServiceServer
}

func (s *echoServer) Echo(
	ctx context.Context,
	req *pb.EchoRequest,
) (*pb.EchoResponse,error){
	log.Println("Server revieved: %s", req.Message)
	return &pb.EchoResponse{
		Message: "You said" + req.Message,
	},nil
}

func (s *echoServer)UpperCase(
	ctx context.Context,
	req *pb.EchoRequest,
) (*pb.EchoResponse, error){
	
	return &pb.EchoResponse{
		Message : strings.ToUpper(req.Message),
	}, nil
}


func main(){
	lis, err := net.Listen("tcp", ":50051")
	if err != nil{
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterEchoServiceServer(grpcServer, &echoServer{})

	log.Println("Echo sever running on :50051")
	grpcServer.Serve(lis)
}