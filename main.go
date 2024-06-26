package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/givek/intro-to-microservices/handlers"
	"github.com/gorilla/mux"
)

func main() {

	logger := log.New(os.Stdout, "product-api ", log.LstdFlags)

	productsHandler := handlers.NewProducts(logger)

	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", productsHandler.GetProducts)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", productsHandler.AddProduct)
	postRouter.Use(productsHandler.ProductValidationMiddleware)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", productsHandler.UpdateProduct)
	putRouter.Use(productsHandler.ProductValidationMiddleware)

	server := &http.Server{
		Addr:         "127.0.0.1:9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {

		err := server.ListenAndServe()

		if err != nil {

			logger.Fatal(err)
		}

	}()

	sigChan := make(chan os.Signal)

	// signal.Notify will broadcast a message on the provided channel,
	// when the provided signal occurs.
	signal.Notify(sigChan, os.Interrupt)

	signal.Notify(sigChan, os.Kill)

	// reading from a channel will block util a message is available to be consumed.
	sig := <-sigChan

	logger.Println("Received terminate, graceful shutdown", sig)

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.Shutdown(timeoutContext)

	// Basic way to create a http server is to use this fucntion.
	// It takes two args:
	//	1. An address string
	//	2. handler (http handler).
	// http.ListenAndServe("127.0.0.1:9090", serveMux)

}
