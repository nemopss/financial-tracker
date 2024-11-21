package main

import (
	"log"
	"net/http"

	"github.com/nemopss/financial-tracker/config"
)

func main() {
	cfg := config.LoadConfig()

	http.HandleFunc("/ping", pingHandler)

	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong\n"))
}
