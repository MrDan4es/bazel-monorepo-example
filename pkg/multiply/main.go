package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	pb "github.com/mrdan4es/bazel-monorepo-example/api/services/multiply/v1"
)

const port = ":22222"

type MultiplyServer struct {
	pb.UnimplementedMultiplyServiceServer
}

func (*MultiplyServer) Multiply(_ context.Context, r *pb.MultiplyRequest) (*pb.MultiplyResponse, error) {
	return &pb.MultiplyResponse{Result: int64(r.First * r.Second)}, nil
}

func NewMultiplyServer() *MultiplyServer {
	return &MultiplyServer{}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("create listener: %v", err)
	}

	serverRegister := grpc.NewServer()

	go func() {
		<-ctx.Done()
		log.Println("gracefully stopping gRPC server")
		serverRegister.GracefulStop()
	}()

	s := NewMultiplyServer()
	pb.RegisterMultiplyServiceServer(serverRegister, s)

	log.Printf("gRPC server started on port %s", port)
	if err = serverRegister.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Fatal(err)
	}
}
