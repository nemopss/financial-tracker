package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nemopss/financial-tracker/internal/repository"
)

type AnalyticsHandler struct {
	Repo repository.Repository
}

// Handlers

// GetIncomeAndExpensesGin handles fetching income and expenses analytics.
// @Summary Get income and expenses
// @Description Fetch total income and expenses for the authenticated user
// @Tags Analytics
// @Produce json
// @Security BearerAuth
// @Success 200 {object} repository.Analytics
// @Router /analytics/income-expenses [get]
func (h *AnalyticsHandler) GetIncomeAndExpensesGin(c *gin.Context) {
	userID := c.GetInt("userID")

	analytics, err := h.Repo.GetIncomeAndExpenses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetIncomeAndExpensesFilteredGin handles fetching income and expenses analytics within a date range.
// @Summary Get income and expenses (filtered)
// @Description Fetch total income and expenses within a specific date range for the authenticated user
// @Tags Analytics
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "Start date in YYYY-MM-DD format"
// @Param end_date query string true "End date in YYYY-MM-DD format"
// @Success 200 {object} repository.Analytics
// @Router /analytics/income-expenses-filtered [get]
func (h *AnalyticsHandler) GetIncomeAndExpensesFilteredGin(c *gin.Context) {
	userID := c.GetInt("userID")

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date and end date are required"})
		return
	}

	analytics, err := h.Repo.GetIncomeAndExpensesFiltered(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch filtered analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetCategoryAnalyticsGin handles fetching category-based analytics.
// @Summary Get category analytics
// @Description Fetch total expenses per category for the authenticated user
// @Tags Analytics
// @Produce json
// @Security BearerAuth
// @Success 200 {array} repository.CategoryAnalytics
// @Router /analytics/categories [get]
func (h *AnalyticsHandler) GetCategoryAnalyticsGin(c *gin.Context) {
	userID := c.GetInt("userID")

	analytics, err := h.Repo.GetCategoryAnalytics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetCategoryAnalyticsFilteredGin handles fetching category-based analytics within a date range.
// @Summary Get category analytics (filtered)
// @Description Fetch total expenses per category within a specific date range for the authenticated user
// @Tags Analytics
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "Start date in YYYY-MM-DD format"
// @Param end_date query string true "End date in YYYY-MM-DD format"
// @Success 200 {array} repository.CategoryAnalytics
// @Router /analytics/categories-filtered [get]
func (h *AnalyticsHandler) GetCategoryAnalyticsFilteredGin(c *gin.Context) {
	userID := c.GetInt("userID")

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date and end date are required"})
		return
	}

	analytics, err := h.Repo.GetCategoryAnalyticsFiltered(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}
