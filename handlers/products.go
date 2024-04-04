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

	p.logger.Println("Handle GET Products.")

	productsList := data.GetProducts()

	err := productsList.ToJSON(rw)

	// jsData, err := json.Marshal(productsList)

	if err != nil {
		http.Error(rw, "Failed to parse JSON!", http.StatusInternalServerError)
		return
	}

	// rw.Write(jsData)

}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {

	p.logger.Println("Handle POST Product.")

	product := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&product)

}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {

	p.logger.Println("Handle PUT Product.")

	vars := mux.Vars(r)

	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)

	if err != nil {

		p.logger.Fatalf("Failed to convert id str to int. id: %v :: err: %v", idStr, err)

		http.Error(rw, "Product ID not found.", http.StatusBadRequest)
		return
	}

	product := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &product)

	if err != nil {

		p.logger.Fatalf("Failed update product. id: %v :: product: %v :: err: %v", id, product, err)

		http.Error(rw, "Failed to update Product", http.StatusBadRequest)
		return
	}

}

type KeyProduct struct{}

func (p *Products) ProductValidationMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		product := data.Product{}

		err := product.FromJson(r.Body)

		if err != nil {

			p.logger.Fatalf("Failed to parse JSON. req body: %v :: err: %v", r.Body, err)

			http.Error(rw, "Failed to parse JSON.", http.StatusBadRequest)

			return

		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)

		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)

	})

}
