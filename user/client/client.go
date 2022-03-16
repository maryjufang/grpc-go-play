package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/maryjufang/grpc-go-play/user/userpb"
	"google.golang.org/grpc"
)

const (
	address = "0.0.0.0:50051"
)

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := userpb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	params := &userpb.GetRequest{}
	r, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("could not get users: %v", err)
	}
	log.Print("\nUSER LIST: \n")
	fmt.Printf("r.GetUsers(): %v\n", r.GetUsers())
}
