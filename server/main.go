package main

import (
	"errors"

	//"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"context"
	"grpc-echo/server/models"
	"log"
	"net"
	"strings"

	pb "grpc-echo/proto"

	"google.golang.org/grpc"
)

type echoServer struct{
	pb.UnimplementedEchoServiceServer
	db *gorm.DB
}

func (s *echoServer) Echo(
	ctx context.Context,
	req *pb.EchoRequest,
) (*pb.EchoResponse,error){
	msg := models.Message{
		Content: req.Message,
	}
	if err := s.db.Create(&msg).Error ; err != nil{
		return nil, err
	}
	return &pb.EchoResponse{
		Message: "Saved" + req.Message,
	}, nil


}

func (s *echoServer)UpperCase(
	ctx context.Context,
	req *pb.EchoRequest,
) (*pb.EchoResponse, error){
	
	return &pb.EchoResponse{
		Message : strings.ToUpper(req.Message),
	}, nil
}

func (s *echoServer) GetMessages(
	ctx context.Context,
	req *pb.GetMessagesRequest,
) (*pb.GetMessagesResponse, error) {

	limit := int(req.Limit)
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	var msgs []models.Message
	if err := s.db.Order("id desc").Limit(limit).Find(&msgs).Error; err != nil {
		return nil, err
	}

	out := make([]*pb.StoredMessage, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, &pb.StoredMessage{
			Id:          uint32(m.ID),
			Content:     m.Content,
			CreatedUnix: m.CreatedAt.Unix(),
		})
	}

	return &pb.GetMessagesResponse{Messages: out}, nil
}

func (s *echoServer) MarkMessageAsRead(
	ctx context.Context,
	req *pb.MarkMessageAsReadRequest,
) (*pb.MarkMessageAsReadResponse, error){
	var msg models.Message
	err := s.db.First(&msg, req.MessageId).Error
	
	if err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil, status.Error(codes.NotFound, "message not found")
		}
		return nil, status.Error(codes.Internal, "failed to load message")
	}
	msg.Read = true
	if err := s.db.Save(&msg).Error; err != nil{
		return nil, status.Error(codes.Internal, "failed to update message")
	}
	return &pb.MarkMessageAsReadResponse{
		Messages: &pb.StoredMessage{
			Id : uint32(msg.ID),
			Content: msg.Content,
			CreatedUnix: msg.CreatedAt.Unix(),
			Read: msg.Read,
		},

	}, nil
}

func(s *echoServer) GetMessageByID(
	ctx  context.Context,
	req *pb.GetMessageByIDRequest,
)(*pb.GetMessageByIDResponse, error){
	var msg models.Message
	err := s.db.First(&msg, req.MessageId).Error
	
	if err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil ,status.Error(codes.NotFound, "failed to load message")
		}
		return nil, status.Error(codes.Internal, "failed to load message")
	}
	return &pb.GetMessageByIDResponse{
		Message: &pb.StoredMessage{
			Id : uint32(msg.ID),
			Content: msg.Content,
			CreatedUnix: msg.CreatedAt.Unix(),
			Read: msg.Read,
		},
	}, nil


}
	




func main(){
	db, err := gorm.Open(sqlite.Open("messages.db"), &gorm.Config{})
	
	if err != nil{
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Message{})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil{
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterEchoServiceServer(grpcServer, &echoServer{db: db})

	log.Println("Echo sever running on :50051")
	go startGateway()
	grpcServer.Serve(lis)

	


}