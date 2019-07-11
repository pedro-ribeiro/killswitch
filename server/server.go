package main

import (
	restapi "killswitch/rest_api"
	"killswitch/store"
	"log"
)

func main() {
	failed := make(chan bool, 1)
	store, err := store.NewRedisStore("featurestore", "localhost:6379")

	if err != nil {
		log.Fatalf("Could not initialize feature store: %s", err)
		panic(err)
	}

	//REST
	restapi.BindAPI("8080", store, failed)
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
