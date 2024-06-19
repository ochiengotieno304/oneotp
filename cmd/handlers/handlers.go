package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ochiengotieno304/oneotp/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CustomMatcher(key string) (string, bool) {
	switch key {
	case "client_d":
		return key, true
	case "secret_key":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func RunHandlers() {
	conn, err := grpc.NewClient(
		"0.0.0.0:6000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	ctx := context.Background()

	gwmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(CustomMatcher),
	)

	// REGISTER GATEWAYS
	err = pb.RegisterAccountServiceHandler(ctx, gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register otp gateway:", err)
	}
	err = pb.RegisterOTPServiceHandler(ctx,gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register otp gateway:", err)
	}
	
	gwServer := &http.Server{
		Addr:    ":6090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on port 6090")
	log.Fatalln(gwServer.ListenAndServe())
}
