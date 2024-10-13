package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	currencyProtos "github.com/givek/intro-to-microservices/currency-api/protos/currency/protos"
	"github.com/givek/intro-to-microservices/products-api/data"
	"github.com/givek/intro-to-microservices/products-api/handlers"
	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	grpcConn, err := grpc.NewClient(
		"localhost:9092",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		panic(err)
	}
	defer grpcConn.Close()

	currencyClient := currencyProtos.NewCurrencyClient(grpcConn)

	productsDB := data.NewProductsDB(currencyClient, logger)

	productsHandler := handlers.NewProducts(logger, productsDB)

	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()

	getRouter.HandleFunc("/", productsHandler.GetProducts).Queries("currency", "[A-Z]{3}")
	getRouter.HandleFunc("/", productsHandler.GetProducts)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()

	postRouter.HandleFunc("/", productsHandler.AddProduct)
	postRouter.Use(productsHandler.ProductValidationMiddleware)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()

	putRouter.HandleFunc("/{id:[0-9]+}", productsHandler.UpdateProduct)
	putRouter.Use(productsHandler.ProductValidationMiddleware)

	// CORS
	corsHandler := gorillahandlers.CORS(
		gorillahandlers.AllowedOrigins([]string{"http://localhost:5173"}),
	)

	server := &http.Server{
		Addr:         "127.0.0.1:9000",
		Handler:      corsHandler(serveMux),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		ErrorLog:     logger,
	}

	go func() {

		err := server.ListenAndServe()

		if err != nil {
			logger.Fatal(err)
		}

	}()

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	logger.Println("Received terminate signal, shuting down gracefully.", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	server.Shutdown(tc)

}
