package main

import (
	"log"
	"net/http"
)

var datastore = NewDataStore()

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "PUT" {
		webUpload(w, req)
		return
	}

	if req.Method == "GET" {
		webDownload(w, req)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

func main() {
	http.HandleFunc("/data/", handler)
	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":8282", nil))
}
