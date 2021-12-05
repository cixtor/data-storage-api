package main

import (
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	http.HandleFunc("/data/", handler)
	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":8282", nil))
}
