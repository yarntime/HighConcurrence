package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/yarntime/HighConcurrence/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/yarntime/HighConcurrence/pkg"
	"github.com/yarntime/HighConcurrence/pkg/jobs"
)

var (
	port       = flag.Int("port", 50051, "The server port")
)

var dispatch *pkg.Dispatcher

type server struct{}

func (s *server) GetInfo(ctx context.Context, in *pb.InfoRequest) (*pb.InfoResponse, error) {
	log.Printf("Client: %s", in.Name)
	job := jobs.NewListJob(in.Name)

	dispatch.Submit(job)

	result := <- job.ResultChan
	return &pb.InfoResponse{Version: "1.2", Endpoint: result.(string)}, nil
}

func main() {

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	dispatch = pkg.NewDispatcher()
	dispatch.Run()

	var opts []grpc.ServerOption

	s := grpc.NewServer(opts...)
	pb.RegisterInfoServer(s, &server{})
	s.Serve(lis)
}
