package main

import (
	"log"
	"time"
	"math/rand"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/aklyukin/rabbit-test-proto"
)

const (
	address     = "localhost:50051"
	nodeName = "first-node"
)

func registerNode(client pb.ServerClient, nodeName *pb.RegisterNodeRequest) bool{
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

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewServerClient(conn)
	//cServer := pb.

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		time.Sleep(time.Duration(r.Int63n(10)) * time.Second)
		isRegistered := registerNode(client, &pb.RegisterNodeRequest{nodeName})
		if isRegistered == true {
			break
		}
	}




	//
	//stream, err := client.PingNode(context.Background()
	//waitc := make(chan struct{})
	//
	//c := pb.NewClientClient(conn)
	//
	//// Contact the server and print out its response.
	//name := defaultName
	//if len(os.Args) > 1 {
	//	name = os.Args[1]
	//}
	//r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	//if err != nil {
	//	log.Fatalf("could not greet: %v", err)
	//}
	//log.Printf("Greeting: %s", r.Message)
}