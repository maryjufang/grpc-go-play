package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/maryjufang/grpc-go-play/user/userpb"
	"google.golang.org/grpc"
)

const (
	address = "0.0.0.0:50051"
)

func argParser(n1 string, n2 string, n3 string) (string, string, int32) {
	first_name := os.Args[1]
	last_name := os.Args[2]
	age, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("Cannot parse arge[3]: %s", err)
	}
	return first_name, last_name, int32(age)
}

func showUsers(c userpb.UserManagementClient) {
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

func doCreate(c userpb.UserManagementClient, first_name, last_name string, age int32) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &userpb.CreateRequest{
		User: &userpb.User{
			FirstName: first_name,
			LastName:  last_name,
			Age:       age,
		},
	}
	r, err := c.CreateUser(ctx, req)
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}
	log.Printf("Create response: %v", r.User.FirstName)
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("2 names and 1 number expected: firstName, lastName, age")
	}
	first_name, last_name, age := argParser(os.Args[1], os.Args[2], os.Args[3])
	log.Printf("first_name = %v, last_name = %v, age = %v ", first_name, last_name, age)

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := userpb.NewUserManagementClient(conn)

	showUsers(c)
	doCreate(c, first_name, last_name, age)
	log.Printf("After Create")
	showUsers(c)
}
