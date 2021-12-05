package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type UploadResponse struct {
	OID  ObjectID `json:"oid"`
	Size int      `json:"size"`
}

func webUpload(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 3 || parts[2] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repository := RepositoryID(parts[2])

	// Limit ioutil to read a maximum of 20MiB of data.
	reader := io.LimitReader(r.Body, 2<<20)
	data, err := ioutil.ReadAll(reader)

	if err != nil {
		// Something weird happened with the IO reader.
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("webUpload", "ioutil.ReadAll", err)
		return
	}

	oid := datastore.Create(repository, data)
	out := UploadResponse{
		OID:  oid,
		Size: len(data),
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.Println("webUpload", "json.Encode", err)
		return
	}
}
