package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/givek/intro-to-microservices/files-api/files"
	"github.com/givek/intro-to-microservices/files-api/handlers"
	gorillahandlers "github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
)

var basePath = env.String("BASE_PATH", false, "./imagestore", "Base path to save images")

func main() {

	env.Parse()

	logger := log.New(os.Stdout, "files-api", log.LstdFlags)

	serveMux := mux.NewRouter()

	serveMux.UseEncodedPath()

	store, err := files.NewLocal(*basePath, 1024*1000*5)

	if err != nil {
		logger.Fatal("Unable to create local file stoare")
		return
	}

	fileHandler := handlers.NewFiles(store, logger)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fileHandler.ServeHTTP)

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", http.StripPrefix("/images/", http.FileServer(http.Dir(*basePath))).ServeHTTP)

	// CORS
	corsHandler := gorillahandlers.CORS(
		gorillahandlers.AllowedOrigins([]string{"http://localhost:5173"}),
	)

	server := &http.Server{
		Addr:         "127.0.0.1:9090",
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
