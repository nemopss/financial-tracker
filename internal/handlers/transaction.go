package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/nemopss/financial-tracker/internal/middleware"
	"github.com/nemopss/financial-tracker/internal/repository"
	"github.com/nemopss/financial-tracker/internal/response"
)

type TransactionHandler struct {
	Repo *repository.DB
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var txn repository.Transaction
	if err := json.NewDecoder(r.Body).Decode(&txn); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	txn.UserID = userID
	txn.Date = time.Now()

	id, err := h.Repo.CreateTransaction(context.Background(), txn)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	response.Success(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	transactions, err := h.Repo.GetTransactions(context.Background(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	response.Success(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var txn repository.Transaction
	if err := json.NewDecoder(r.Body).Decode(&txn); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	txn.UserID = userID
	txn.Date = time.Now()

	if err := h.Repo.UpdateTransaction(context.Background(), txn); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update transaction")
		return
	}

	response.Success(w, http.StatusNoContent, nil)
}

func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	txnID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	if err := h.Repo.DeleteTransaction(context.Background(), userID, txnID); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete transaction")
		return
	}

	response.Success(w, http.StatusNoContent, nil)
}
