package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"sync"
	"time"

	pbMaster "github.com/aklyukin/rabbit-test-proto/master"

	"github.com/aklyukin/rabbit-test/int/bidi"
)

const (
	port = ":50051"
	serverName = "MainServer"
)

// server is used to implement Server.
type server struct{
	pbMaster.UnimplementedMasterServer
}

// RegisterNode implements Server
func (s *server) RegisterNode(ctx context.Context, in *pbMaster.RegisterNodeRequest) (*pbMaster.RegisterNodeResponse, error) {
	log.Printf("[GRPC] Client ask for register node: %s", in.NodeName)
	//curTime := time.Now().Unix()
	//if curTime % 3 == 0{
	//	log.Printf("[GRPC] Node registered: %s", in.NodeName)
	//	return &pbMaster.RegisterNodeResponse{IsRegistered: true}, nil
	//	//	add nodename to list
	//}
	//return &pbMaster.RegisterNodeResponse{}, nil
	return &pbMaster.RegisterNodeResponse{IsRegistered: true}, nil
}

func (s *server) PingServer(ctx context.Context, in *pbMaster.Ping) (*pbMaster.Pong, error){
	log.Printf("[GRPC] Ping from client: %s", in.NodeName)
	return &pbMaster.Pong{ServerName: serverName}, nil
}

func main() {

	grpcServer := grpc.NewServer()
	pbMaster.RegisterMasterServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	var nodesMap = &sync.Map{}

	gconn := bidi.Listen(port, nodesMap, grpcServer)
	defer gconn.Close()

	//grpcClient := pbNode.NewNodeClient(gconn)

	log.Printf("Server started")

	// Waiting for node[s]
	time.Sleep(10 * time.Second)
	for {
		time.Sleep(10 * time.Second)
		log.Printf("Nodes: ")
		nodesMap.Range(func(key interface{}, value interface{}) bool {
			log.Printf(fmt.Sprint(key))
			return true
		})
	}

	time.Sleep(10 * time.Minute)
}

//https://github.com/chetan/bidi-hello/tree/master/helloworld