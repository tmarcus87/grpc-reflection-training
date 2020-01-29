package main

import (
	"context"
	helloworld "github.com/tmarcus87/grpc-reflection-training/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type GreeterServer struct {
}

func (g GreeterServer) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{ Message: "Hello, " + in.Name + "!" }, nil
}

func main() {
	listen, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}

	ctrl := GreeterServer{}

	server := grpc.NewServer()
	helloworld.RegisterGreeterServer(server, &ctrl)
	reflection.Register(server)
	if err := server.Serve(listen); err != nil {
		panic(err)
	}
}
