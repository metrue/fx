package main

import (
	"context"
	"log"
	"time"

	"github.com/metrue/fx/api"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8866"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewAPIClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ListService(ctx, &api.ListServiceRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %v", r.GetServices())
}
