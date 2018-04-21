package main

import (
	"log"
	"net"
	"time"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/aklyukin/rabbit-test-proto"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) RegisterNode(ctx context.Context, in *pb.RegisterNodeRequest) (*pb.RegisterNodeReply, error) {
	log.Printf("Client ask for register node: %s", in.NodeName)
	curTime := time.Now().Unix()
	if curTime % 3 == 0{
		log.Printf("register node")
		return &pb.RegisterNodeReply{true}, nil
		//	add nodename to list
	}
	return &pb.RegisterNodeReply{false}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServerServer(s, &server{})
	log.Printf("Server started")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}