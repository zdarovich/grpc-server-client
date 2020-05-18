package main

import (
	"fmt"
	"github.com/zdarovich/grpc-server-client/internal/api"
	"github.com/zdarovich/grpc-server-client/internal/server"
	"google.golang.org/grpc"
	"log"
	"net"
)


func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	s := server.Server{}
	// create a gRPC server object
	grpcServer := grpc.NewServer()
	api.RegisterProxyServer(grpcServer, &s)
	// start the server
	log.Print("server started")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
