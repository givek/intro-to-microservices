package server

import (
	"context"
	"log"

	protos "github.com/givek/intro-to-microservices/currency-api/protos/currency/protos"
)

type Currency struct {
	logger *log.Logger
	protos.UnimplementedCurrencyServer
}

func NewCurrency(logger *log.Logger) *Currency {
	return &Currency{logger: logger}
}

func (c *Currency) GetRate(_ context.Context, reqRate *protos.RateRequest) (*protos.RateResponse, error) {

	c.logger.Println("GetRate", reqRate.GetBase(), reqRate.GetDestination())

	var f float32 = 0.30

	return &protos.RateResponse{Rate: f}, nil

}
