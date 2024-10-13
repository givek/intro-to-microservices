package data

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type ExchangeRates struct {
	logger *log.Logger
	rates  map[string]float64 // TODO: check if we can use Currency enum instead of string.
}

func (e *ExchangeRates) GetRate(base, dest string) (float64, error) {
	br, ok := e.rates[base]

	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", base)
	}

	dr, ok := e.rates[dest]

	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", dest)
	}

	return dr / br, nil

}
func NewExchangeRates(logger *log.Logger) (*ExchangeRates, error) {

	exchangeRates := &ExchangeRates{logger: logger, rates: map[string]float64{}}

	err := exchangeRates.getRates()

	if err != nil {
		return nil, err
	}

	return exchangeRates, nil

}

func (e *ExchangeRates) getRates() error {
	ratesRes, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")

	if err != nil {
		return err
	}

	if ratesRes.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code 200 got %d", ratesRes.StatusCode)
	}

	defer ratesRes.Body.Close()

	md := &Cubes{}
	xml.NewDecoder(ratesRes.Body).Decode(&md)

	for _, c := range md.CubeData {

		r, err := strconv.ParseFloat(c.Rate, 64)

		if err != nil {

			return err

		}

		e.rates[c.Currency] = r

	}

	e.rates["EUR"] = 1

	return nil

}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}
