package handlers

import (
	"log"
	"net/http"

	"github.com/givek/intro-to-microservices/data"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	return &Products{logger: logger}
}

func (p *Products) getProducts(rw http.ResponseWriter) {

	productsList := data.GetProducts()

	err := productsList.ToJSON(rw)

	// jsData, err := json.Marshal(productsList)

	if err != nil {
		http.Error(rw, "Failed to parse JSON!", http.StatusInternalServerError)
		return
	}

	// rw.Write(jsData)

}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		p.getProducts(rw)
		return

	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)

	}

}
