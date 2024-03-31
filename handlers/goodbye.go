package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type Goodbye struct {
	logger *log.Logger
}

func NewGoodbye(logger *log.Logger) *Goodbye {

	return &Goodbye{logger: logger}

}

func (g *Goodbye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	g.logger.Println("Goodbye handler called!")

	fmt.Fprintf(rw, "Goodbye")

}
