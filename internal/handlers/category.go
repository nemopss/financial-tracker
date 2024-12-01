package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nemopss/financial-tracker/internal/repository"
)

type CategoryHandler struct {
	Repo repository.Repository
}

// Models for Swagger documentation

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required" example:"Groceries"`
}

type CreateCategoryResponse struct {
	ID int `json:"id" example:"1"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required" example:"Updated Category"`
}

type CategoryListResponse struct {
	Categories []repository.Category `json:"categories"`
}

// Handlers

// CreateCategoryGin handles category creation using Gin framework.
// @Summary Create a new category
// @Description Create a category for the authenticated user
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body CreateCategoryRequest true "Category data"
// @Success 201 {object} CreateCategoryResponse
// @Router /categories [post]
func (h *CategoryHandler) CreateCategoryGin(c *gin.Context) {
	userID := c.GetInt("userID") // Получаем userID из middleware

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	id, err := h.Repo.CreateCategory(c.Request.Context(), userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, CreateCategoryResponse{ID: id})
}

// GetCategoriesGin handles fetching all categories for the user.
// @Summary Get categories
// @Description Fetch all categories for the authenticated user
// @Tags Categories
// @Produce json
// @Security BearerAuth
// @Success 200 {object} CategoryListResponse
// @Router /categories/list [get]
func (h *CategoryHandler) GetCategoriesGin(c *gin.Context) {
	userID := c.GetInt("userID")

	categories, err := h.Repo.GetCategories(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, CategoryListResponse{Categories: categories})
}

// UpdateCategoryGin handles updating a category for the user.
// @Summary Update a category
// @Description Update a category by ID for the authenticated user
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "Category ID"
// @Param category body UpdateCategoryRequest true "Category data"
// @Success 204 "No Content"
// @Router /categories/update [put]
func (h *CategoryHandler) UpdateCategoryGin(c *gin.Context) {
	userID := c.GetInt("userID")

	categoryID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.Repo.UpdateCategory(c.Request.Context(), userID, categoryID, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteCategoryGin handles deleting a category for the user.
// @Summary Delete a category
// @Description Delete a category by ID for the authenticated user
// @Tags Categories
// @Produce json
// @Security BearerAuth
// @Param id query int true "Category ID"
// @Success 204 "No Content"
// @Router /categories/delete [delete]
func (h *CategoryHandler) DeleteCategoryGin(c *gin.Context) {
	userID := c.GetInt("userID")

	categoryID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.Repo.DeleteCategory(c.Request.Context(), userID, categoryID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.Status(http.StatusNoContent)
}

