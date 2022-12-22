package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"grpc_go/pbs/pb"
	"grpc_go/service"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func loadcredentials() (credentials.TransportCredentials, error) {
	servercert, err := tls.LoadX509KeyPair("/cert/server-cert.pem", "/cert/server-key.pem")
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{servercert},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(config), nil
}
func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("starting the server on the port %d", *port)
	laptopserver := service.NewLaptop()
	tlscertificate, err := loadcredentials()
	if err != nil {
		log.Fatal("cannot load credentials", err)
	}
	grpcserver := grpc.NewServer(grpc.Creds(tlscertificate))
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
