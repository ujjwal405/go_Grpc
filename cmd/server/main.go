package main

import (
	"flag"
	"fmt"
	"grpc_go/pbs/pb"
	"grpc_go/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("starting the server on the port %d", *port)
	laptopserver := service.NewLaptop()
	grpcserver := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcserver, laptopserver)
	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server")
	}
	err = grpcserver.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
