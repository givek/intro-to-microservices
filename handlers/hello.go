package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct{ logger *log.Logger }

func NewHello(logger *log.Logger) *Hello {

	return &Hello{logger: logger}

}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	d, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(rw, "Unknow Error!", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(rw, "Hello %s", d)
}
