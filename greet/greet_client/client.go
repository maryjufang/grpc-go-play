package main

import (
	"context"
	"io"
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

func doServerStreaming(c greetpb.GreetServiceClient) {
	log.Printf("In Client doServerStreaming")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mary",
			LastName:  "Chen",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
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

	doServerStreaming(c)
}
