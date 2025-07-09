package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func startHTTPServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	log.Println("🌐 HTTP server started at http://localhost:8092")
	log.Fatal(http.ListenAndServe(":8092", nil))
}
