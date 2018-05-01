package main

import (
	"log"
	"time"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/aklyukin/rabbit-test-proto"

	"github.com/aklyukin/rabbit-test/int/bidi"
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
	return &pb.Pong{serverName}, nil
}

func main() {

	grpcServer := grpc.NewServer()
	pb.RegisterServerServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	gconn := bidi.Listen(port, grpcServer)
	defer gconn.Close()

	grpcClient := pb.NewNodeClient(gconn)

	log.Printf("Server started")
	// Waiting for node[s]
	time.Sleep(1 * time.Minute)
	for {
		r, err := grpcClient.CheckStatus(context.Background(), &pb.Empty{})
		if err != nil {
			log.Printf("could not check status: %v", err)
		}
		log.Printf("Check node status: %s", r.Status)
		time.Sleep(10 * time.Second)
	}

	time.Sleep(10 * time.Minute)
}