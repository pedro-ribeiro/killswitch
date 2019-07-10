package main

import (
	restapi "killswitch/rest_api"
	"log"
)

func main() {
	failed := make(chan bool, 1)

	//REST
	restapi.BindAPI("8080", failed)
	// err := rest_api.BindAPI("8080")

	//gRPC

	// flag.Parse()
	// lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8777))
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	// grpcServer := grpc.NewServer()
	// pb.RegisterFeaturesServer(grpcServer, &featuresService{})
	// grpcServer.Serve(lis)

	// if err == nil {
	log.Printf("REST API listening on port %s\n", "8080")

	<-failed
	// }

}

// type featuresService struct {
// }

// func (s *featuresService) GetFeatures(request *pb.FeaturesRequest, server features.Features_GetFeaturesServer) error {
// 	return nil
// }
