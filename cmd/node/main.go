package main

import (
	"os"
	"log"
	"time"
	"math/rand"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/aklyukin/rabbit-test-proto"

)

const (
	address     = "localhost:50051"
)

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
	nodeName := "first-node"
	if len(os.Args) > 1 {
		nodeName = os.Args[1]
	} else {
		nodeName = "first-node"
	}

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
				pongServer := PingServer(client, &pb.Ping{"kjkjkjk"})
				log.Printf("Pong from server: %s", pongServer)
			}
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