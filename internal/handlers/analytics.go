package handlers

import (
	"context"
	"net/http"

	"github.com/nemopss/financial-tracker/internal/middleware"
	"github.com/nemopss/financial-tracker/internal/repository"
	"github.com/nemopss/financial-tracker/internal/response"
)

type AnalyticsHandler struct {
	Repo repository.Repository
}

func (h *AnalyticsHandler) GetIncomeAndExpenses(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	analytics, err := h.Repo.GetIncomeAndExpenses(context.Background(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch analytics")
		return
	}

	response.Success(w, http.StatusOK, analytics)
}

func (h *AnalyticsHandler) GetIncomeAndExpensesFiltered(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		response.Error(w, http.StatusBadRequest, "Start date and end date are required")
		return
	}

	analytics, err := h.Repo.GetIncomeAndExpensesFiltered(context.Background(), userID, startDate, endDate)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch filtered analytics")
		return
	}

	response.Success(w, http.StatusOK, analytics)
}

func (h *AnalyticsHandler) GetCategoryAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	analytics, err := h.Repo.GetCategoryAnalytics(context.Background(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch analytics")
		return
	}

	response.Success(w, http.StatusOK, analytics)
}

func (h *AnalyticsHandler) GetCategoryAnalyticsFiltered(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		response.Error(w, http.StatusBadRequest, "Start date and end date are required")
		return
	}

	analytics, err := h.Repo.GetCategoryAnalyticsFiltered(context.Background(), userID, startDate, endDate)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch analytics")
		return
	}

	response.Success(w, http.StatusOK, analytics)
}
