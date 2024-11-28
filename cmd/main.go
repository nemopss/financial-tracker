package main

import (
	"log"
	"net/http"

	"github.com/nemopss/financial-tracker/config"
	"github.com/nemopss/financial-tracker/internal/handlers"
	"github.com/nemopss/financial-tracker/internal/repository"
)

func main() {
	cfg := config.LoadConfig()

	db, err := repository.NewDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Conn.Close()

	authHandler := &handlers.AuthHandler{
		Repo:      db,
		JWTSecret: cfg.JWTSecret,
	}

	http.HandleFunc("/api/v1/register", authHandler.Register)
	http.HandleFunc("/api/v1/login", authHandler.Login)

	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
