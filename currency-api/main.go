package main

import (
	"fmt"
	"log"
	"net"
	"os"

	protos "github.com/givek/intro-to-microservices/currency-api/protos/currency/protos"
	"github.com/givek/intro-to-microservices/currency-api/server"

	"google.golang.org/grpc"
)

func main() {

	logger := log.New(os.Stdout, "currency-api", log.LstdFlags)

	grpcServer := grpc.NewServer()

	currencyServer := server.NewCurrency(logger)

	protos.RegisterCurrencyServer(grpcServer, currencyServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 9092))

	if err != nil {
		logger.Println("Failed to create listener", err)
		return
	}

	grpcServer.Serve(listener)
}
