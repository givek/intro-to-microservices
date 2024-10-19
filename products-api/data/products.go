package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"

	currencyProtos "github.com/givek/intro-to-microservices/currency-api/protos/currency/protos"
	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) FromJson(reader io.Reader) error {

	decoder := json.NewDecoder(reader)

	return decoder.Decode(p)

}

func validateSKU(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)

	matches := regex.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}

func (p *Product) Validate() error {
	v := validator.New()

	v.RegisterValidation("sku", validateSKU)

	return v.Struct(p)
}

type Products []*Product

type ProductsDB struct {
	currencyClient currencyProtos.CurrencyClient
	logger         *log.Logger
	rates          map[string]float64
	client         currencyProtos.Currency_SubscribeRatesClient
}

func NewProductsDB(
	currencyClient currencyProtos.CurrencyClient,
	logger *log.Logger,
) *ProductsDB {
	productsDB := &ProductsDB{
		logger:         logger,
		currencyClient: currencyClient,
		rates:          make(map[string]float64),
		client:         nil,
	}

	go productsDB.handleUpdates()

	return productsDB
}

func (p *Products) ToJson(writer io.Writer) error {

	encoder := json.NewEncoder(writer)

	return encoder.Encode(p)
}

func (p *ProductsDB) handleUpdates() {
	subClient, err := p.currencyClient.SubscribeRates(context.Background())

	if err != nil {

		p.logger.Println("Unable subscribe for rate updates", err)

	}

	p.client = subClient

	for {
		res, err := subClient.Recv()

		// because we used oneof is proto def
		// it is guaranteed that either error or message is populated.

		if res.GetError() != nil {

			p.logger.Println("Received an error message from server", res.GetError())

			// start listening for the next message.
			continue

		}

		if rateRes := res.GetRateResponse(); rateRes != nil {

			p.logger.Println("Received updated rate from server", rateRes.GetDestination().String(), rateRes.Rate)

			if err != nil {
				p.logger.Println("Failed to receive rate response", err)
			}

			p.rates[rateRes.Destination.String()] = float64(rateRes.Rate)
		}
	}
}

func (p *ProductsDB) getRate(dest string) (float32, error) {

	if r, ok := p.rates[dest]; ok {
		return float32(r), nil
	}

	// Get the exchange rate
	rateReq := &currencyProtos.RateRequest{
		Base:        currencyProtos.Currencies_EUR,
		Destination: currencyProtos.Currencies(currencyProtos.Currencies_value[dest]),
	}

	// get initial rate
	rateRes, err := p.currencyClient.GetRate(
		context.Background(),
		rateReq,
	)

	p.rates[rateReq.Destination.String()] = float64(rateRes.Rate)

	// subscribe for updates
	p.client.Send(rateReq)

	if err != nil {
		return 0, err
	}

	return rateRes.Rate, nil

}

func (p *ProductsDB) GetProducts(currency string) (Products, error) {

	if currency == "" {
		return productList, nil
	}

	rate, err := p.getRate(currency)

	if err != nil {
		return nil, err
	}

	productListCopy := []*Product{}

	for _, p := range productList {

		// creates a copy of p
		pCopy := *p

		pCopy.Price = pCopy.Price * rate

		productListCopy = append(productListCopy, &pCopy)
	}

	return productListCopy, nil
}

func AddProduct(p *Product) {

	id := len(productList) + 1

	p.ID = id

	productList = append(productList, p)

}

func UpdateProduct(id int, p *Product) error {

	_, pos, err := findProduct(id)

	if err != nil {
		return err
	}

	p.ID = id

	productList[pos] = p

	return nil

}

var ErrProductNotFound = fmt.Errorf("product not found")

func findProduct(id int) (*Product, int, error) {

	for i, p := range productList {
		if p.ID == id {

			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound

}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffe",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
