package main

import (
	"log"
	"net"
	"time"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/aklyukin/rabbit-test-proto"

	"github.com/aklyukin/rabbit-test/internal/bidi"
)

const (
	port = ":50051"
	serverName = "MainServer"
)

// server is used to implement Server.
type server struct{}

// RegisterNode implements Server
func (s *server) RegisterNode(ctx context.Context, in *pb.RegisterNodeRequest) (*pb.RegisterNodeResponse, error) {
	log.Printf("Client ask for register node: %s", in.NodeName)
	curTime := time.Now().Unix()
	if curTime % 3 == 0{
		log.Printf("register node")
		return &pb.RegisterNodeResponse{true}, nil
		//	add nodename to list
	}
	return &pb.RegisterNodeResponse{false}, nil
}

func (s *server) PingServer(ctx context.Context, in *pb.Ping) (*pb.Pong, error){
	log.Printf("Ping from client: %s", in.NodeName)
	return &pb.Pong{"jjjjjjj"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterServerServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	gconn := bidi.Listen(port, grpcServer)
	defer gconn.Close()


	log.Printf("Server started")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}