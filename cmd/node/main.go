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
	log.Printf("Server ask for node status")
	return &pb.NodeStatusResponse{pb.NodeStatusResponse_READY}, nil
}

func RegisterNode(client pb.ServerClient, nodeName *pb.RegisterNodeRequest) bool{
	resp, err := client.RegisterNode(context.Background(), nodeName)
	if err != nil {
		log.Fatalf("Error registering node: %v", err)
	}
	if resp.IsRegistered == true {
		log.Printf("Node registered")
	}
	if resp.IsRegistered == false {
		log.Printf("Node not registered")
	}
	return resp.IsRegistered
}

func PingServer(client pb.ServerClient, nodeName *pb.Ping) string{
	resp, err := client.PingServer(context.Background(), nodeName)
	if err != nil {
		log.Fatalf("Error ping server: %v", err)
	}
	return resp.ServerName
}

func main() {
	nodeName := "default-node"
	if len(os.Args) > 1 {
		if os.Args[1] == "--plain" {
			doPlainClient(address, nodeName)
			return
		}
		nodeName = os.Args[1]
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNodeServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	// open channel and create client
	gconn := bidi.Connect(address, grpcServer)
	defer gconn.Close()
	grpcClient := pb.NewServerClient(gconn)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		time.Sleep(time.Duration(r.Int63n(10)) * time.Second)
		log.Printf("Try to register node")
		isRegistered := RegisterNode(grpcClient, &pb.RegisterNodeRequest{nodeName})
		if isRegistered == true {
			for {
				time.Sleep(5 * time.Second)
				pongServer := PingServer(grpcClient, &pb.Ping{nodeName})
				log.Printf("Pong from server: %s", pongServer)
			}
		}
	}
}

func doPlainClient(address string, nodeName string) {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewServerClient(conn)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		time.Sleep(time.Duration(r.Int63n(10)) * time.Second)
		isRegistered := RegisterNode(client, &pb.RegisterNodeRequest{nodeName})
		if isRegistered == true {
			for {
				time.Sleep(5 * time.Second)
				pongServer := PingServer(client, &pb.Ping{nodeName})
				log.Printf("Pong from server: %s", pongServer)
			}
		}
	}
}