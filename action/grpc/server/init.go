package server

import (
	pb "hcc/cello/action/grpc/pb/rpccello"
	"hcc/cello/lib/config"
	"hcc/cello/lib/logger"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Init : Initialize gRPC server
func Init() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(int(config.Grpc.Port)))
	if err != nil {
		logger.Logger.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCelloServer(s, &celloServer{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	logger.Logger.Println("Opening gRPC server on port " + strconv.Itoa(int(config.Grpc.Port)) + "...")
	if err := s.Serve(lis); err != nil {
		logger.Logger.Fatalf("failed to serve: %v", err)
	}
}
