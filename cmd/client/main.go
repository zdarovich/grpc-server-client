package main

import (
	"context"
	"flag"
	"github.com/zdarovich/grpc-server-client/internal/api"
	"github.com/zdarovich/grpc-server-client/internal/client"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	url := flag.String("url", "https://www.google.ee", "url sent to server by grpc ")
	flag.Parse()

	go run()
	log.Printf("proxy %s", *url)

	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	// create stream
	cli := api.NewProxyClient(conn)

	_, err = cli.Init(context.Background(), &api.UrlMessage{
		Url: *url,
	})
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}
}

func run() error {

	lis, err := net.Listen("tcp", ":7778")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	c := client.Client{}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	api.RegisterProxyCallerServer(grpcServer, &c)
	// start the server
	log.Print("start listening")

	return grpcServer.Serve(lis)
}
