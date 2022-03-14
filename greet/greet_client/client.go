package main

import (
	"log"

	"github.com/maryjufang/grpc-go-play/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	log.Printf("In Greet Client func")

	cc, err := grpc.Dial("locahost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	log.Printf("Created Client: %f", c)
}
