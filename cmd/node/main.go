package main

import (
	pbMaster "github.com/aklyukin/rabbit-test-proto/master"
	pbNode "github.com/aklyukin/rabbit-test-proto/node"
	"golang.org/x/net/context"
	"log"
	"os"

	"github.com/aklyukin/rabbit-test/int/bidi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"time"
)

const (
	address     = "localhost:50051"
)

type server struct{}

func (s *server) CheckStatus(ctx context.Context, empty *pbNode.Empty) (*pbNode.NodeStatusResponse, error){
	log.Printf("[GRPC] Server ask for node status")
	return &pbNode.NodeStatusResponse{Status: pbNode.NodeStatusResponse_READY}, nil
}

func RegisterNode(client pbMaster.MasterClient, nodeName *pbMaster.RegisterNodeRequest) bool{
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

func PingServer(client pbMaster.MasterClient, nodeName *pbMaster.Ping) string{
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
	pbNode.RegisterNodeServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	// open channel and create client
	gconn := bidi.Connect(address, grpcServer)
	defer gconn.Close()

	grpcClient := pbMaster.NewMasterClient(gconn)

	log.Printf("Try to register node")
	time.Sleep(3 * time.Second)
	isRegistered := RegisterNode(grpcClient, &pbMaster.RegisterNodeRequest{NodeName: nodeName})
	if isRegistered == true {
		for {
			time.Sleep(5 * time.Second)
			PingServer(grpcClient, &pbMaster.Ping{NodeName: nodeName})
		}
	}
}