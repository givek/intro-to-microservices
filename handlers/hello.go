package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	logger *log.Logger
}

func NewHello(logger *log.Logger) *Hello {

	return &Hello{logger: logger}

}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	h.logger.Println("Hello handler called!")

	data, err := io.ReadAll(r.Body) // implements the interface io.ReadCloser. So standard go libararies for reading from that stream.

	if err != nil {
		http.Error(rw, "Opps!", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(rw, "Hello %s\n", data) // printing the string to the ResponseWriter.

}
