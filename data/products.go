package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Product defines the structure for an API product
// ` ` After the type def are call struct tags
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) FromJson(r io.Reader) error {

	jsDecoder := json.NewDecoder(r)

	return jsDecoder.Decode(p)

}

// Products is a collection of Product
type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {

	jsEncoder := json.NewEncoder(w)

	return jsEncoder.Encode(p)

}

// productList is a hard coded list of products for this
// example data source
var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

func GetProducts() Products {

	return productList

}

func getNextID() int {

	maxID := 0

	for _, p := range productList {

		maxID = max(maxID, p.ID)

	}

	return maxID + 1

}

func findProduct(id int) (*Product, error) {

	for _, p := range productList {

		if p.ID == id {

			return p, nil

		}

	}

	return &Product{}, fmt.Errorf("Not found product.")

}

func removeProduct(id int) {

	filteredProducts := []*Product{}

	for _, p := range productList {

		if p.ID != id {

			filteredProducts = append(filteredProducts, p)

		}

	}

	productList = filteredProducts

}

func AddProduct(p *Product) {

	p.ID = getNextID()

	productList = append(productList, p)

}

func UpdateProduct(id int, p *Product) error {

	_, err := findProduct(id)

	if err != nil {
		return err
	}

	removeProduct(id)

	p.ID = id

	productList = append(productList, p)

	return nil

}
