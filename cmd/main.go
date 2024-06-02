package main

import (
	"log"
	"net"

	"github.com/ochiengotieno304/oneotp/cmd/servers"
	"github.com/ochiengotieno304/oneotp/internal/interceptors"
	"github.com/ochiengotieno304/oneotp/internal/utils"
	"github.com/ochiengotieno304/oneotp/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	accountServer := servers.NewAccountServer()
	otpServer := servers.NewOTPServer()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.AuthInterceptor(),
			interceptors.ServerLogInterceptor(utils.InitLogger()),
		),
	)

	reflection.Register(s)

	pb.RegisterAccountServiceServer(s, accountServer)
	pb.RegisterOTPServiceServer(s, otpServer)

	log.Println("Serving gRPC on 0.0.0.0:6000")
	log.Fatalln(s.Serve(lis))
}
