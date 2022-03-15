package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	calculatorpb "github.com/maryjufang/grpc-go-play/calculator/calculatorpb"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func argParser(n1 string, n2 string) (int32, int32) {
	N1, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Cannot parse arge[1]: %s", err)
	}
	N2, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Cannot parse arge[2]: %s", err)
	}
	return int32(N1), int32(N2)
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("2 numbers expected: n1 n2")
	}

	n1, n2 := argParser(os.Args[1], os.Args[2])
	log.Printf("n1 = %v, n2 = %v", n1, n2)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot connect: %s", err)
	}
	defer conn.Close()

	client := calculatorpb.NewCalculatorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reqAdd := &calculatorpb.AddRequest{
		N1: int32(n1),
		N2: int32(n2),
	}

	addResult, err := client.Add(ctx, reqAdd)
	if err != nil {
		log.Fatalf("Adding error: %s", err)
	}
	log.Printf("%d + %d = %d", n1, n2, addResult.N1)

	reqSub := &calculatorpb.SubtractRequest{
		N1: int32(n1),
		N2: int32(n2),
	}

	subtractResult, err := client.Subtract(ctx, reqSub)
	if err != nil {
		log.Fatalf("Subtracting error: %s", err)
	}
	log.Printf("%d - %d = %d", n1, n2, subtractResult.N1)

	reqMul := &calculatorpb.MultiplyRequest{
		N1: int32(n1),
		N2: int32(n2),
	}

	multiplyResult, err := client.Multiply(ctx, reqMul)
	if err != nil {
		log.Fatalf("Multiplying error: %s", err)
	}
	log.Printf("%d * %d = %d", n1, n2, multiplyResult.N1)

	reqDiv := &calculatorpb.DivideRequest{
		N1: int32(n1),
		N2: int32(n2),
	}
	divideResult, err := client.Divide(ctx, reqDiv)
	if err != nil {
		log.Fatalf("Dividing error: %s", err)
	}
	log.Printf("%d / %d = %d", n1, n2, divideResult.N1)

	reqMod := &calculatorpb.ModRequest{
		N1: int32(n1),
		N2: int32(n2),
	}
	modResult, err := client.Mod(ctx, reqMod)
	if err != nil {
		log.Fatalf("Dividing error: %s", err)
	}
	log.Printf("%d %% %d = %d", n1, n2, modResult.N1)
}
