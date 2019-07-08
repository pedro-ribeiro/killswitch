package main

import (
	"flag"
	"fmt"

	// "killswitch/features"
	"killswitch/features"
	pb "killswitch/features"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterFeaturesServer(grpcServer, &featuresService{})
	grpcServer.Serve(lis)
}

type featuresService struct {
}

func (s *featuresService) GetFeatures(request *pb.FeaturesRequest, server features.Features_GetFeaturesServer) error {
	return nil
}
