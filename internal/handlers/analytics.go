package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nemopss/financial-tracker/internal/middleware"
	"github.com/nemopss/financial-tracker/internal/repository"
)

type AnalyticsHandler struct {
	Repo *repository.DB
}

func (h *AnalyticsHandler) GetIncomeAndExpenses(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	analytics, err := h.Repo.GetIncomeAndExpenses(context.Background(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch analytics", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(analytics)
}

func (h *AnalyticsHandler) GetCategoryAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	analytics, err := h.Repo.GetCategoryAnalytics(context.Background(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch analytics", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(analytics)
}
