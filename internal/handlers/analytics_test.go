package handlers

import (
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

func TestGetIncomeAndExpensesHandler(t *testing.T) {
	mockRepo := &repository.MockRepo{}

	mockRepo.On("GetIncomeAndExpenses", mock.Anything, 1).
		Return(&repository.Analytics{TotalIncome: 1000, TotalExpense: -500}, nil)

	handler := &AnalyticsHandler{Repo: mockRepo}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/income-expenses", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))

	w := httptest.NewRecorder()
	handler.GetIncomeAndExpenses(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var wrappedResponse struct {
		StatusCode int                  `json:"status_code"`
		Data       repository.Analytics `json:"data"`
	}
	err := json.NewDecoder(resp.Body).Decode(&wrappedResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Equal(t, 1000.0, wrappedResponse.Data.TotalIncome)
	assert.Equal(t, -500.0, wrappedResponse.Data.TotalExpense)

	mockRepo.AssertExpectations(t)
}

func TestGetIncomeAndExpensesFilteredHandler(t *testing.T) {
	mockRepo := &repository.MockRepo{}

	mockRepo.On("GetIncomeAndExpensesFiltered", mock.Anything, 1, "2024-01-01", "2024-12-31").
		Return(&repository.Analytics{TotalIncome: 5000, TotalExpense: -2500}, nil)

	handler := &AnalyticsHandler{Repo: mockRepo}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/income-expenses-filtered?start_date=2024-01-01&end_date=2024-12-31", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))

	w := httptest.NewRecorder()
	handler.GetIncomeAndExpensesFiltered(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var wrappedResponse struct {
		StatusCode int                  `json:"status_code"`
		Data       repository.Analytics `json:"data"`
	}
	err := json.NewDecoder(resp.Body).Decode(&wrappedResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Equal(t, 5000.0, wrappedResponse.Data.TotalIncome)
	assert.Equal(t, -2500.0, wrappedResponse.Data.TotalExpense)

	mockRepo.AssertExpectations(t)
}

func TestGetCategoryAnalyticsHandler(t *testing.T) {
	mockRepo := &repository.MockRepo{}

	mockAnalytics := []repository.CategoryAnalytics{
		{CategoryName: "Groceries", TotalAmount: -200.50},
		{CategoryName: "Entertainment", TotalAmount: -100.00},
	}

	mockRepo.On("GetCategoryAnalytics", mock.Anything, 1).
		Return(mockAnalytics, nil)

	handler := &AnalyticsHandler{Repo: mockRepo}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/categories", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))

	w := httptest.NewRecorder()
	handler.GetCategoryAnalytics(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var wrappedResponse struct {
		StatusCode int                            `json:"status_code"`
		Data       []repository.CategoryAnalytics `json:"data"`
	}
	err := json.NewDecoder(resp.Body).Decode(&wrappedResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Len(t, wrappedResponse.Data, 2)
	assert.Equal(t, "Groceries", wrappedResponse.Data[0].CategoryName)
	assert.Equal(t, -200.50, wrappedResponse.Data[0].TotalAmount)
	assert.Equal(t, "Entertainment", wrappedResponse.Data[1].CategoryName)
	assert.Equal(t, -100.00, wrappedResponse.Data[1].TotalAmount)

	mockRepo.AssertExpectations(t)
}

func TestGetCategoryAnalyticsFilteredHandler(t *testing.T) {
	mockRepo := &repository.MockRepo{}

	mockAnalytics := []repository.CategoryAnalytics{
		{CategoryName: "Groceries", TotalAmount: -500.00},
		{CategoryName: "Entertainment", TotalAmount: -200.00},
	}

	mockRepo.On("GetCategoryAnalyticsFiltered", mock.Anything, 1, "2024-01-01", "2024-12-31").
		Return(mockAnalytics, nil)

	handler := &AnalyticsHandler{Repo: mockRepo}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/categories-filtered?start_date=2024-01-01&end_date=2024-12-31", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))

	w := httptest.NewRecorder()
	handler.GetCategoryAnalyticsFiltered(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var wrappedResponse struct {
		StatusCode int                            `json:"status_code"`
		Data       []repository.CategoryAnalytics `json:"data"`
	}
	err := json.NewDecoder(resp.Body).Decode(&wrappedResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Len(t, wrappedResponse.Data, 2)
	assert.Equal(t, "Groceries", wrappedResponse.Data[0].CategoryName)
	assert.Equal(t, -500.00, wrappedResponse.Data[0].TotalAmount)
	assert.Equal(t, "Entertainment", wrappedResponse.Data[1].CategoryName)
	assert.Equal(t, -200.00, wrappedResponse.Data[1].TotalAmount)

	mockRepo.AssertExpectations(t)
}
