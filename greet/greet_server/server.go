package main

import (
	"log"
	"net"

	"github.com/maryjufang/grpc-go-play/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	log.Printf("In Greet Server func")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}
}
