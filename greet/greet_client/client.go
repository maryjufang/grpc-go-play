package main

import (
	"context"
	"log"

	"github.com/maryjufang/grpc-go-play/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func doGreet(c greetpb.GreetServiceClient) {
	log.Printf("In Client doGreet")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mary",
			LastName:  "Chen",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Greet failed: %v", err)
	}
	log.Printf("Greet Response: %v", res.Result)
}

func main() {
	log.Printf("In Greet Client func")

	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)

	doGreet(c)
}
