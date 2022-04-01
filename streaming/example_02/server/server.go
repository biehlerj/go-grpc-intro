package main

import (
	"fmt"
	"io"
	"net"

	pb "github.com/biehlerj/go-grpc-intro/streaming/example_02/numbers"
	"google.golang.org/grpc"
)

type NumServer struct {
	pb.UnimplementedNumServiceServer
}

func (n *NumServer) Sum(stream pb.NumService_SumServer) error {
	var total int64 = 0
	var counter int64 = 0
	for {
		next, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("Received %d numbers sum: %d\n", counter, total)
			stream.SendAndClose(&pb.NumResponse{Total: total})
			return nil
		}
		if err != nil {
			return err
		}
		total += next.X
		counter++
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()

	pb.RegisterNumServiceServer(s, &NumServer{})

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
