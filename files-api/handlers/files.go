package handlers

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

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

func (f *Files) saveFile(id, path string, rw http.ResponseWriter, rc io.ReadCloser) {

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, rc)

	if err != nil {
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}

func (f *Files) UploadREST(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	fileName := vars["filename"]

	f.logger.Println("Upload file REST", id, fileName)

	// No need to check for invalid id or filename as the mux router will not send requests
	// here unless they have the correct parameters.

	f.saveFile(id, fileName, rw, r.Body)
}

func (f *Files) UploadMultipart(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(128 * 1024)

	if err != nil {
		http.Error(rw, "Failed to parse the multipart form", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(rw, "Failed to parse the provided id", http.StatusBadRequest)
		return
	}

	f.logger.Println("Process for id", id)

	uploadedFile, multipartFileHeader, err := r.FormFile("file")

	if err != nil {
		http.Error(rw, "Please upload a file.", http.StatusBadRequest)
		return
	}

	f.saveFile(strconv.Itoa(id), multipartFileHeader.Filename, rw, uploadedFile)

}
