package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/givek/intro-to-microservices/data"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	return &Products{logger: logger}
}

func (p *Products) getProducts(rw http.ResponseWriter) {

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

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {

	p.logger.Println("Handle POST Product.")

	product := &data.Product{}

	err := product.FromJson(r.Body)

	if err != nil {
		http.Error(rw, "Failed to parse JSON.", http.StatusBadRequest)
		return
	}

	data.AddProduct(product)

}

func (p *Products) updateProduct(rw http.ResponseWriter, r *http.Request) {

	p.logger.Println("Handle PUT Product.")

	regx := regexp.MustCompile(`/[0-9]+`)

	group := regx.FindAllStringSubmatch(r.URL.Path, -1)

	if len(group) != 1 {

		p.logger.Printf("Invalid URL, more than 1 ids found. group: %v", group)

		http.Error(rw, "Please send exactly 1 product ID.", http.StatusBadRequest)
		return
	}

	if len(group[0][0]) != 2 {

		p.logger.Printf("Invalid URL, more than 2 captures. group[0]: %v", group[0])

		http.Error(rw, "Product ID not found.", http.StatusBadRequest)
		return
	}

	idStr := string(group[0][0][1])

	id, err := strconv.Atoi(idStr)

	if err != nil {

		p.logger.Fatalf("Failed to convert id str to int. id: %v :: err: %v", idStr, err)

		http.Error(rw, "Product ID not found.", http.StatusBadRequest)
		return
	}

	product := &data.Product{}

	err = product.FromJson(r.Body)

	if err != nil {

		p.logger.Fatalf("Failed to parse JSON. req body: %v :: err: %v", r.Body, err)

		http.Error(rw, "Failed to parse JSON.", http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(id, product)

	if err != nil {

		p.logger.Fatalf("Failed update product. id: %v :: product: %v :: err: %v", id, product, err)

		http.Error(rw, "Failed to update Product", http.StatusBadRequest)
		return
	}

}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		p.getProducts(rw)
		return

	case http.MethodPost:
		p.addProduct(rw, r)
		return

	case http.MethodPut:
		p.updateProduct(rw, r)
		return

	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return

	}

}
