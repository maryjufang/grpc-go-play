package main

import (
	"context"
	"errors"
	"log"
	"net"

	calculatorpb "github.com/maryjufang/grpc-go-play/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	protocol = "tcp"
	port     = ":50051"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServer
}

func (s *server) Add(ctx context.Context, in *calculatorpb.AddRequest) (*calculatorpb.AddReply, error) {
	log.Printf("Add n1 = %v, n2 = %v", in.N1, in.N2)
	return &calculatorpb.AddReply{N1: in.N1 + in.N2}, nil
}

func (s *server) Subtract(ctx context.Context, in *calculatorpb.SubtractRequest) (*calculatorpb.SubtractReply, error) {
	log.Printf("Subtract n1 = %v, n2 = %v", in.N1, in.N2)
	return &calculatorpb.SubtractReply{N1: in.N1 - in.N2}, nil
}

func (s *server) Multiply(ctx context.Context, in *calculatorpb.MultiplyRequest) (*calculatorpb.MultiplyReply, error) {
	log.Printf("Multiply n1 = %v, n2 = %v", in.N1, in.N2)
	return &calculatorpb.MultiplyReply{N1: in.N1 * in.N2}, nil
}

func (s *server) Divide(ctx context.Context, in *calculatorpb.DivideRequest) (*calculatorpb.DivideReply, error) {
	log.Printf("Divide n1 = %v, n2 = %v", in.N1, in.N2)
	if in.N2 == 0 {
		return &calculatorpb.DivideReply{N1: 0}, errors.New("ERROR: cannot divide by zero.")
	}
	return &calculatorpb.DivideReply{N1: in.N1 / in.N2}, nil
}

func (s *server) Mod(ctx context.Context, in *calculatorpb.ModRequest) (*calculatorpb.ModReply, error) {
	log.Printf("Modulo n1 = %v, n2 = %v", in.N1, in.N2)
	if in.N2 == 0 {
		return &calculatorpb.ModReply{N1: 0}, errors.New("ERROR: cannot mod by zero.")
	}
	return &calculatorpb.ModReply{N1: in.N1 % in.N2}, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.Calculator_PrimeNumberDecompositionServer) error {
	log.Printf("Received PrimeNumberDecomposition RPC: %v\n", req)

	number := req.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			log.Printf("Divisor has increased to %v\n", divisor)
		}
	}
	return nil
}

func main() {
	log.Printf("Calculator server, waiting for inputs")
	lis, err := net.Listen(protocol, port)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
