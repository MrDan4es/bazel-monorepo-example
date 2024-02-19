package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	pb "github.com/mrdan4es/bazel-monorepo-example/api/services/minus/v1"
)

const port = ":33333"

type MinusServer struct {
	pb.UnimplementedMinusServiceServer
}

func (*MinusServer) Minus(_ context.Context, r *pb.MinusRequest) (*pb.MinusResponse, error) {
	return &pb.MinusResponse{Result: int64(r.First - r.Second)}, nil
}

func NewMinusServer() *MinusServer {
	return &MinusServer{}
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

	s := NewMinusServer()
	pb.RegisterMinusServiceServer(serverRegister, s)

	log.Printf("gRPC server started on port %s", port)
	if err = serverRegister.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Fatal(err)
	}
}
