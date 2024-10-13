// Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
// 	BasePath: /
// 	Version: 1.0.0
//
// 	Consumes:
//	- application/json
//
// 	Produces:
//	- application/json
//
// swagger:meta

package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/givek/intro-to-microservices/products-api/data"

	"github.com/gorilla/mux"
)

type Products struct {
	logger     *log.Logger
	productsDB *data.ProductsDB
}

func NewProducts(
	logger *log.Logger,
	productsDB *data.ProductsDB,
) *Products {
	return &Products{
		logger:     logger,
		productsDB: productsDB,
	}
}

// swagger:route GET /products products listProducts
// Returns a list of products

// GetProducts returns the products from the data store.
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {

	p.logger.Println("Get Products")

	currency := r.URL.Query().Get("currency")

	products, err := p.productsDB.GetProducts(currency)

	if err != nil {
		p.logger.Println("Failed to fetch products from DB", err)
		http.Error(rw, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	err = products.ToJson(rw)

	if err != nil {
		p.logger.Println("Failed to convert products to JSON", err)
		http.Error(rw, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {

	p.logger.Println("Post Product")

	product := r.Context().Value(KeyProduct{}).(data.Product)

	// product := &data.Product{}

	// err := product.FromJson(r.Body)

	// if err != nil {
	// 	http.Error(rw, "Unable to parse the request body.", http.StatusBadRequest)
	// 	return
	// }

	data.AddProduct(&product)

}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {

	p.logger.Println("PUT Product")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Invalid Id.", http.StatusBadRequest)
		return
	}

	product := r.Context().Value(KeyProduct{}).(data.Product)

	// product := &data.Product{}

	// err = product.FromJson(r.Body)

	// if err != nil {
	// 	http.Error(rw, "Unable to parse the request body.", http.StatusBadRequest)
	// 	return
	// }

	err = data.UpdateProduct(id, &product)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Unable to update product", http.StatusInternalServerError)
		return
	}

}

type KeyProduct struct{}

func (p *Products) ProductValidationMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		product := data.Product{}

		err := product.FromJson(r.Body)

		if err != nil {
			http.Error(rw, "Unable to parse the request body.", http.StatusBadRequest)
			return
		}

		err = product.Validate()

		if err != nil {
			http.Error(
				rw,
				fmt.Sprintf("Failed to validate request body. %s", err),
				http.StatusBadRequest,
			)

			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)

		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)

	})

}
