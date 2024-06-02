package main

import (
	"log"
	"net"

	"github.com/ochiengotieno304/oneotp/cmd/servers"
	"github.com/ochiengotieno304/oneotp/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	accountServer := servers.NewAccountServer()
	authServer := servers.NewAuthServer()

	s := grpc.NewServer()
	reflection.Register(s)

	pb.RegisterAccountServiceServer(s, accountServer)
	pb.RegisterAuthServiceServer(s, authServer)

	log.Println("Serving gRPC on 0.0.0.0:6000")
	log.Fatalln(s.Serve(lis))
}
