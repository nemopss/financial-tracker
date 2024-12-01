package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nemopss/financial-tracker/internal/middleware"
	"github.com/nemopss/financial-tracker/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionHandler(t *testing.T) {
	mockRepo := &repository.MockRepo{}

	// Проверка транзакции без точного сравнения даты
	mockRepo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(txn repository.Transaction) bool {
		return txn.Amount == 100.50 &&
			txn.Description == "Groceries" &&
			txn.CategoryID == 1 &&
			txn.UserID == 1
	})).Return(1, nil)

	handler := &TransactionHandler{Repo: mockRepo}

	body := map[string]interface{}{
		"amount":      100.50,
		"description": "Groceries",
		"category_id": 1,
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))

	w := httptest.NewRecorder()
	handler.CreateTransaction(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var wrappedResponse struct {
		StatusCode int
		Data       map[string]int
	}
	err := json.NewDecoder(resp.Body).Decode(&wrappedResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Equal(t, 1, wrappedResponse.Data["id"])

	mockRepo.AssertExpectations(t)
}

func TestGetTransactionsHandler(t *testing.T) {
	mockRepo := &repository.MockRepo{}

	mockTransactions := []repository.Transaction{
		{ID: 1, Amount: 100.50, Description: "Groceries", CategoryID: 1, UserID: 1},
		{ID: 2, Amount: -15.32, Description: "Coffee", CategoryID: 2, UserID: 1},
	}
	mockRepo.On("GetTransactions", mock.Anything, 1).
		Return(mockTransactions, nil)

	handler := &TransactionHandler{Repo: mockRepo}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/list", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))

	w := httptest.NewRecorder()
	handler.GetTransactions(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var wrappedResponse struct {
		StatusCode int                      `json:"status_code"`
		Data       []repository.Transaction `json:"data"`
	}
	err := json.NewDecoder(resp.Body).Decode(&wrappedResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Len(t, wrappedResponse.Data, 2)
	assert.Equal(t, "Groceries", wrappedResponse.Data[0].Description)
	assert.Equal(t, "Coffee", wrappedResponse.Data[1].Description)

	mockRepo.AssertExpectations(t)
}
