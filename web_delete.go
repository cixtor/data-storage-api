package main

import (
	"errors"
	"net/http"
	"strings"
)

// webDelete performs an idempotent delete operation against
// the repository map and its corresponding object map.
func webDelete(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 4 || parts[2] == "" || parts[3] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repository := RepositoryID(parts[2])
	objectID := ObjectID(parts[3])

	err := datastore.Delete(repository, objectID)

	if errors.Is(err, errRepositoryNotFound) || errors.Is(err, errObjectNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
