package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nemopss/financial-tracker/internal/repository"
)

// Models for Swagger documentation

type TransactionHandler struct {
	Repo repository.Repository
}

type CreateTransactionRequest struct {
	Amount      float64 `json:"amount" example:"100.50"`
	Date        string  `json:"date,omitempty" example:"2024-12-01T15:04:05Z"`
	Description string  `json:"description" example:"Grocery shopping"`
	CategoryID  int     `json:"category_id" example:"1"`
}

type CreateTransactionResponse struct {
	ID int `json:"id" example:"1"`
}

// Handlers

// CreateTransactionGin handles the creation of a transaction.
// @Summary Create a new transaction
// @Description Add a new transaction for the authenticated user
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction body CreateTransactionRequest true "Transaction data (date is optional, defaults to now)"
// @Success 201 {object} CreateTransactionResponse
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransactionGin(c *gin.Context) {
	userID := c.GetInt("userID")

	var txn repository.Transaction
	if err := c.ShouldBindJSON(&txn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	txn.UserID = userID
	txn.Date = time.Now()

	id, err := h.Repo.CreateTransaction(c.Request.Context(), txn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetTransactionsGin handles fetching all transactions for the user.
// @Summary Get transactions
// @Description Fetch all transactions for the authenticated user
// @Tags Transactions
// @Produce json
// @Security BearerAuth
// @Success 200 {array} repository.Transaction
// @Router /transactions/list [get]
func (h *TransactionHandler) GetTransactionsGin(c *gin.Context) {
	userID := c.GetInt("userID")

	transactions, err := h.Repo.GetTransactions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// UpdateTransactionGin handles updating a transaction.
// @Summary Update a transaction
// @Description Update a transaction for the authenticated user
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction body repository.Transaction true "Transaction data (include transaction ID)"
// @Success 204 "No Content"
// @Router /transactions/update [put]
func (h *TransactionHandler) UpdateTransactionGin(c *gin.Context) {
	userID := c.GetInt("userID")

	var txn repository.Transaction
	if err := c.ShouldBindJSON(&txn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	txn.UserID = userID
	txn.Date = time.Now()

	if err := h.Repo.UpdateTransaction(c.Request.Context(), txn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteTransactionGin handles deleting a transaction.
// @Summary Delete a transaction
// @Description Delete a transaction by ID for the authenticated user
// @Tags Transactions
// @Produce json
// @Security BearerAuth
// @Param id query int true "Transaction ID"
// @Success 204 "No Content"
// @Router /transactions/delete [delete]
func (h *TransactionHandler) DeleteTransactionGin(c *gin.Context) {
	userID := c.GetInt("userID")

	txnID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	if err := h.Repo.DeleteTransaction(c.Request.Context(), userID, txnID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	c.Status(http.StatusNoContent)
}
