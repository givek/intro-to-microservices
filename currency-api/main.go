package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/givek/intro-to-microservices/currency-api/data"
	protos "github.com/givek/intro-to-microservices/currency-api/protos/currency/protos"
	"github.com/givek/intro-to-microservices/currency-api/server"

	"google.golang.org/grpc"
)

func main() {

	logger := log.New(os.Stdout, "currency-api", log.LstdFlags)

	grpcServer := grpc.NewServer()

	exchangeRates, err := data.NewExchangeRates(logger)

	if err != nil {
		logger.Println("Failed to get exchange rates", err)
		return
	}

	currencyServer := server.NewCurrency(logger, exchangeRates)

	protos.RegisterCurrencyServer(grpcServer, currencyServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 9092))

	if err != nil {
		logger.Println("Failed to create listener", err)
		return
	}

	grpcServer.Serve(listener)
}
