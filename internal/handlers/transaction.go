package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/nemopss/financial-tracker/internal/middleware"
	"github.com/nemopss/financial-tracker/internal/repository"
)

type TransactionHandler struct {
	Repo *repository.DB
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var txn repository.Transaction
	if err := json.NewDecoder(r.Body).Decode(&txn); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	txn.UserID = userID
	txn.Date = time.Now()

	id, err := h.Repo.CreateTransaction(context.Background(), txn)
	if err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	transactions, err := h.Repo.GetTransactions(context.Background(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var txn repository.Transaction
	if err := json.NewDecoder(r.Body).Decode(&txn); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	txn.UserID = userID
	txn.Date = time.Now()

	if err := h.Repo.UpdateTransaction(context.Background(), txn); err != nil {
		http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	txnID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.DeleteTransaction(context.Background(), userID, txnID); err != nil {
		http.Error(w, "Failed to delete transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
