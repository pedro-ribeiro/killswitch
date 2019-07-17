package main

import (
	"flag"
	grpcapi "killswitch/grpc_api"
	restapi "killswitch/rest_api"
	"killswitch/store"
	"log"
)

func main() {
	flag.Parse()

	failed := make(chan bool, 1)
	store, err := store.NewRedisStore("featurestore", "localhost:6379")

	if err != nil {
		log.Fatalf("Could not initialize feature store: %s", err)
		panic(err)
	}

	//REST
	restapi.BindAPI("8080", store, failed)

	//gRPC
	grpcapi.BindAPI("8081", store, failed)

	log.Printf("gRPC API listening on port %s\n", "8081")
	log.Printf("REST API listening on port %s\n", "8080")

	<-failed
}
