package main

import (
	"os"
	"log"
	"golang.org/x/net/context"
	pb "github.com/aklyukin/rabbit-test-proto"

	"math/rand"
	"time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/aklyukin/rabbit-test/int/bidi"
)

const (
	address     = "localhost:50051"
)

type server struct{}

func (s *server) CheckStatus(ctx context.Context, empty *pb.Empty) (*pb.NodeStatusResponse, error){
	log.Printf("[GRPC] Server ask for node status")
	return &pb.NodeStatusResponse{pb.NodeStatusResponse_READY}, nil
}

func RegisterNode(client pb.ServerClient, nodeName *pb.RegisterNodeRequest) bool{
	resp, err := client.RegisterNode(context.Background(), nodeName)
	if err != nil {
		log.Fatalf("[GRPC] Error registering node: %v", err)
	}
	if resp.IsRegistered == true {
		log.Printf("[GRPC] Node registered")
	}
	if resp.IsRegistered == false {
		log.Printf("[GRPC] Node not registered")
	}
	return resp.IsRegistered
}

func PingServer(client pb.ServerClient, nodeName *pb.Ping) string{
	resp, err := client.PingServer(context.Background(), nodeName)
	if err != nil {
		log.Fatalf("[GRPC] Error ping server: %v", err)
	}
	log.Printf("[GRPC] Pong from server: %s", resp.ServerName)
	return resp.ServerName
}

func main() {
	nodeName := "default-node"
	if len(os.Args) > 1 {
		nodeName = os.Args[1]
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNodeServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	// open channel and create client
	gconn := bidi.Connect(address, grpcServer)
	defer gconn.Close()
	log.Printf("Test 10")
	grpcClient := pb.NewServerClient(gconn)
	log.Printf("Test 11")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		time.Sleep(time.Duration(r.Int63n(10)) * time.Second)
		log.Printf("Try to register node")
		isRegistered := RegisterNode(grpcClient, &pb.RegisterNodeRequest{nodeName})
		if isRegistered == true {
			for {
				time.Sleep(5 * time.Second)
				PingServer(grpcClient, &pb.Ping{nodeName})
			}
		}
	}
}