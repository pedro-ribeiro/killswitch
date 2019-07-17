package grpcapi

import (
	fmt "fmt"
	"killswitch/features"
	"log"
	"net"

	grpc "google.golang.org/grpc"
)

type grpcServer struct {
	store features.FeatureStore
}

func BindAPI(port string, store features.FeatureStore, failed chan bool) {
	log.Println("Starting GRPC API")

	server := grpcServer{store}

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			failed <- true
		}
		grpcServer := grpc.NewServer()
		RegisterFeaturesServer(grpcServer, &server)
		grpcServer.Serve(lis)
	}()
}

func (s *grpcServer) GetFeatures(request *GrpcFeaturesRequest, stream Features_GetFeaturesServer) error {
	values, err := getFeaturesStream(s.store)

	if err != nil {
		return err
	}

	for f := range values {
		if err := stream.Send(messageFromFeature(f)); err != nil {
			return err
		}
	}

	return nil
}

func messageFromFeature(value features.Feature) *GrpcFeaturesResponse {
	return &GrpcFeaturesResponse{
		Key:         value.Key,
		Description: value.Description,
		IsActive:    value.IsActive,
	}
}

func getFeaturesStream(store features.FeatureStore) (chan features.Feature, error) {
	stream := make(chan features.Feature)
	go func() {
		values, err := store.GetAllFeatures()

		if err != nil {
			fmt.Printf("could not get current features: %s", err)
			close(stream)
			return
		}

		for _, f := range values {
			stream <- f
		}

		close(stream)
	}()
	return stream, nil
}
