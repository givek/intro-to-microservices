package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/givek/intro-to-microservices/data"
	"github.com/gorilla/mux"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	return &Products{logger: logger}
}

func (p *Products) GetProducts(rw http.ResponseWriter, _ *http.Request) {

	p.logger.Println("Get Products")

	products := data.GetProducts()

	err := products.ToJson(rw)

	if err != nil {
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

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)

		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)

	})

}
