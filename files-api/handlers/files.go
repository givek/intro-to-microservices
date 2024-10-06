package handlers

import (
	"log"

	"github.com/givek/intro-to-microservices/files-api/files"
)

type Files struct {
	store  files.Storage
	logger *log.Logger
}

func NewFiles(store files.Storage, logger *log.Logger) *Files {
	return &Files{store: store, logger: logger}
}
