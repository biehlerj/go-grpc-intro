package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	pb "github.com/biehlerj/go-grpc-intro/transcoding/example_01/user"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type UserServer struct {
	httpAddr string
	grpcAddr string
	pb.UnimplementedUserServiceServer
}

func (u *UserServer) Get(ctx context.Context, req *pb.UserRequest) (*pb.User, error) {
	fmt.Println("Server received:", req.String())
	return &pb.User{UserId: "John", Email: "john@gmail.com"}, nil
}

func (u *UserServer) Create(ctx context.Context, req *pb.User) (*pb.User, error) {
	fmt.Println("Server received:", req.String())
	return &pb.User{UserId: req.UserId, Email: req.Email}, nil
}

func (u *UserServer) ServeGrpc() {
	lis, err := net.Listen("tcp", u.grpcAddr)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, u)
	fmt.Println("Server listening GRPC:")

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func (u *UserServer) ServeHttp() {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	endpoint := u.grpcAddr
	err := pb.RegisterUserServiceHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
	if err != nil {
		panic(err)
	}

	httpServer := &http.Server{
		Addr:    u.httpAddr,
		Handler: mux,
	}

	fmt.Println("Server listening HTTP:")
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

func main() {
	us := UserServer{httpAddr: ":8080", grpcAddr: ":50051"}
	go us.ServeGrpc()
	us.ServeHttp()
}
