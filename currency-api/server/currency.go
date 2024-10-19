package server

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/givek/intro-to-microservices/currency-api/data"
	protos "github.com/givek/intro-to-microservices/currency-api/protos/currency/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Currency struct {
	logger *log.Logger
	protos.UnimplementedCurrencyServer

	exchangeRates *data.ExchangeRates

	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

func NewCurrency(
	logger *log.Logger,
	exchangeRates *data.ExchangeRates,
) *Currency {

	c := &Currency{
		logger:        logger,
		exchangeRates: exchangeRates,
		subscriptions: make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest),
	}

	go c.handleUpdates()

	return c
}

func (c *Currency) handleUpdates() {
	ru := c.exchangeRates.MonitorRates(5 * time.Second)

	for range ru {

		c.logger.Println("got updated rate.")

		// loop over subscribed clients
		for k, v := range c.subscriptions {

			// loop over subscribed rates
			for _, rr := range v {

				r, err := c.exchangeRates.GetRate(rr.GetBase().String(), rr.GetDestination().String())

				if err != nil {

					c.logger.Println("Failed to GetRates", err)

				}

				err = k.Send(&protos.StreamingRateResponse{
					Message: &protos.StreamingRateResponse_RateResponse{
						RateResponse: &protos.RateResponse{
							Base:        rr.Base,
							Destination: rr.Destination,
							Rate:        float32(r),
						},
					},
				})

				if err != nil {
					c.logger.Println("Unable to send updated rate", err)
				}

			}

		}

	}

}

func (c *Currency) GetRate(
	_ context.Context,
	reqRate *protos.RateRequest,
) (*protos.RateResponse, error) {

	c.logger.Println("GetRate", reqRate.GetBase(), reqRate.GetDestination())

	if reqRate.Base == reqRate.Destination {

		err := status.Errorf(
			codes.InvalidArgument,
			"Base currency cannot be the same as the Destination currency",
		)

		return nil, err

	}

	var f, err = c.exchangeRates.GetRate(
		reqRate.Base.String(),
		reqRate.Destination.String(),
	)

	if err != nil {
		return nil, err
	}

	rateRes := &protos.RateResponse{
		Rate:        float32(f), // TODO: Not good - float64 to float32
		Base:        reqRate.GetBase(),
		Destination: reqRate.GetDestination(),
	}

	return rateRes, nil

}

func (c *Currency) SubscribeRates(
	src protos.Currency_SubscribeRatesServer,
) error {

	for {

		rr, err := src.Recv()

		if err == io.EOF {

			c.logger.Println("client has closed the connection.")

			break

		}

		if err != nil {

			c.logger.Println("unable to read from client")

			return err

		}

		c.logger.Println("handle client request", rr)

		rrc, ok := c.subscriptions[src]

		if !ok {
			rrc = []*protos.RateRequest{}
		}

		var validationError *status.Status

		// check that subscription does not exists
		for _, v := range rrc {

			if v.Base == rr.Base && v.Destination == rr.Destination {

				c.logger.Println("This subscription already exits.", v)

				validationError = status.Newf(
					codes.AlreadyExists,
					"unable to subscribe for currency as subscription already exists",
				)

				break

			}

		}

		if validationError != nil {

			c.logger.Println("sending validation error.", validationError)

			src.Send(&protos.StreamingRateResponse{
				Message: &protos.StreamingRateResponse_Error{
					Error: validationError.Proto(), // convert the status to status
				},
			})

			continue
		}

		rrc = append(rrc, rr)

		c.subscriptions[src] = rrc
	}

	return nil

}
