package handlers

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/givek/intro-to-microservices/files-api/files"
	"github.com/gorilla/mux"
)

type Files struct {
	store  files.Storage
	logger *log.Logger
}

func NewFiles(store files.Storage, logger *log.Logger) *Files {
	return &Files{store: store, logger: logger}
}

func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r *http.Request) {

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r.Body)

	if err != nil {
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}

func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	fileName := vars["filename"]

	// No need to check for invalid id or filename as the mux router will not send requests
	// here unless they have the correct parameters.

	f.saveFile(id, fileName, rw, r)
}
