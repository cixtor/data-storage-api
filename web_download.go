package main

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

func webDownload(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 4 || parts[2] == "" || parts[3] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repository := RepositoryID(parts[2])
	objectID := ObjectID(parts[3])

	data, err := datastore.Read(repository, objectID)

	if errors.Is(err, errRepositoryNotFound) || errors.Is(err, errObjectNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(data); err != nil {
		log.Println("webDownload", "w.Write", err)
	}
}
