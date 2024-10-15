package server

import (
	"context"
	"io"
	"log"
	"time"

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

func (c *Currency) GetRate(
	_ context.Context,
	reqRate *protos.RateRequest,
) (*protos.RateResponse, error) {

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

func (c *Currency) SubscribeRates(
	src protos.Currency_SubscribeRatesServer,
) error {

	go func() {

		for {

			rr, err := src.Recv()

			if err == io.EOF {

				c.logger.Println("client has closed the connection.")

				break

			}

			if err != nil {

				c.logger.Println("unable to read from client")

				break

			}

			c.logger.Println("handle client request", rr)

		}

	}()

	for {

		err := src.Send(&protos.RateResponse{Rate: 2})

		if err != nil {

			return err

		}

		time.Sleep(5 * time.Second)

	}
}
