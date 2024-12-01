package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nemopss/financial-tracker/internal/middleware"
	"github.com/nemopss/financial-tracker/internal/repository"
	"github.com/nemopss/financial-tracker/internal/response"
)

type CategoryHandler struct {
	Repo repository.Repository
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	id, err := h.Repo.CreateCategory(context.Background(), userID, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	response.Success(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	categories, err := h.Repo.GetCategories(context.Background(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	response.Success(w, http.StatusOK, categories)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	categoryID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.Repo.UpdateCategory(context.Background(), userID, categoryID, req.Name); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update category")
		return
	}

	response.Success(w, http.StatusNoContent, nil)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	categoryID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid categoty ID")
		return
	}

	if err := h.Repo.DeleteCategory(context.Background(), userID, categoryID); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	response.Success(w, http.StatusNoContent, nil)
}
