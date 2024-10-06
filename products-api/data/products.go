package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

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

func (p *Products) ToJson(writer io.Writer) error {

	encoder := json.NewEncoder(writer)

	return encoder.Encode(p)
}

func GetProducts() Products {
	return productList
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

var ErrProductNotFound = fmt.Errorf("Product Not Found.")

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
