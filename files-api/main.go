package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	logger := log.New(os.Stdout, "files-api", log.LstdFlags)

	serveMux := mux.NewRouter()

	postRouter := serveMux.Get(http.MethodPost).Subrouter()
	postRouter.HandleFunc("images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}")

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
