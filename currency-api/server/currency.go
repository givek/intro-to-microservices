package server

import (
	"context"
	"log"

	"github.com/givek/intro-to-microservices/currency-api/data"
	protos "github.com/givek/intro-to-microservices/currency-api/protos/currency/protos"
)

type Currency struct {
	logger *log.Logger
	protos.UnimplementedCurrencyServer

	exchangeRates *data.ExchangeRates
}

func NewCurrency(
	logger *log.Logger,
	exchangeRates *data.ExchangeRates,
) *Currency {
	return &Currency{
		logger:        logger,
		exchangeRates: exchangeRates,
	}
}

func (c *Currency) GetRate(_ context.Context, reqRate *protos.RateRequest) (*protos.RateResponse, error) {

	c.logger.Println("GetRate", reqRate.GetBase(), reqRate.GetDestination())

	var f, err = c.exchangeRates.GetRate(
		reqRate.Base.String(),
		reqRate.Destination.String(),
	)

	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Rate: float32(f)}, nil // TODO: Not good - float64 to float32

}
