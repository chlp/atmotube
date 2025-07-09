package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func startHTTPServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mu.RLock()
		_ = json.NewEncoder(w).Encode(data)
		mu.RUnlock()
	})

	log.Println("🌐 HTTP server started at http://localhost:8092")
	log.Fatal(http.ListenAndServe(":8092", nil))
}
