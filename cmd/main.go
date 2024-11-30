package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/nemopss/financial-tracker/config"
	"github.com/nemopss/financial-tracker/internal/handlers"
	"github.com/nemopss/financial-tracker/internal/middleware"
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

	categoryHandler := &handlers.CategoryHandler{Repo: db}
	transactionHandler := &handlers.TransactionHandler{Repo: db}

	protected := middleware.AuthMiddleware(cfg.JWTSecret)

	// AUTH
	http.HandleFunc("/api/v1/register", authHandler.Register)
	http.HandleFunc("/api/v1/login", authHandler.Login)

	// Categories
	http.Handle("/api/v1/categories", protected(http.HandlerFunc(categoryHandler.CreateCategory)))
	http.Handle("/api/v1/categories/list", protected(http.HandlerFunc(categoryHandler.GetCategories)))
	http.Handle("/api/v1/categories/update", protected(http.HandlerFunc(categoryHandler.UpdateCategory)))
	http.Handle("/api/v1/categories/delete", protected(http.HandlerFunc(categoryHandler.DeleteCategory)))

	// Transactions
	http.Handle("/api/v1/transactions", protected(http.HandlerFunc(transactionHandler.CreateTransaction)))
	http.Handle("/api/v1/transactions/list", protected(http.HandlerFunc(transactionHandler.GetTransactions)))
	http.Handle("/api/v1/transactions/update", protected(http.HandlerFunc(transactionHandler.UpdateTransaction)))
	http.Handle("/api/v1/transactions/delete", protected(http.HandlerFunc(transactionHandler.DeleteTransaction)))

	http.Handle("/api/v1/protected", protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)
		w.Write([]byte("Protected content for user ID: " + strconv.Itoa(userID) + "\n"))
	})))

	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
