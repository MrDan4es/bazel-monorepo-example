package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	pb "github.com/mrdan4es/bazel-monorepo-example/api/services/sum/v1"
)

const port = ":11111"

type SumServer struct {
	pb.UnimplementedSumServiceServer
}

func (*SumServer) Sum(_ context.Context, r *pb.SumRequest) (*pb.SumResponse, error) {
	return &pb.SumResponse{Result: int64(r.First + r.Second)}, nil
}

func NewSumServer() *SumServer {
	return &SumServer{}
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

	s := NewSumServer()
	pb.RegisterSumServiceServer(serverRegister, s)

	log.Printf("gRPC server started on port %s", port)
	if err = serverRegister.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Fatal(err)
	}
}
