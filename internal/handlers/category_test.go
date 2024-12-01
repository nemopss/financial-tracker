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

func TestCreateCategoryHandler(t *testing.T) {
	mockRepo := &repository.MockRepo{}

	mockRepo.On("CreateCategory", mock.Anything, 1, "Groceries").
		Return(1, nil)

	handler := &CategoryHandler{Repo: mockRepo}

	body := map[string]string{"name": "Groceries"}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))

	w := httptest.NewRecorder()
	handler.CreateCategory(w, req)

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

func TestGetCategoriesHandler(t *testing.T) {
	// Создаём mock-репозиторий
	mockRepo := &repository.MockRepo{}

	// Настраиваем mock
	mockCategories := []repository.Category{
		{ID: 1, Name: "Groceries"},
		{ID: 2, Name: "Entertainment"},
	}
	mockRepo.On("GetCategories", mock.Anything, 1).
		Return(mockCategories, nil)

	// Создаём обработчик с mock-репозиторием
	handler := &CategoryHandler{Repo: mockRepo}

	// Создаём запрос
	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories/list", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))

	// Эмулируем HTTP-запрос
	w := httptest.NewRecorder()
	handler.GetCategories(w, req)

	// Проверяем ответ
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Декодируем обёрнутый ответ
	var wrappedResponse struct {
		StatusCode int                   `json:"status_code"`
		Data       []repository.Category `json:"data"`
	}
	err := json.NewDecoder(resp.Body).Decode(&wrappedResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Проверяем содержимое "data"
	assert.Len(t, wrappedResponse.Data, 2)
	assert.Equal(t, "Groceries", wrappedResponse.Data[0].Name)
	assert.Equal(t, "Entertainment", wrappedResponse.Data[1].Name)

	// Проверяем, что mock-методы были вызваны
	mockRepo.AssertExpectations(t)
}

func TestUpdateCategoryHandler(t *testing.T) {
	mockRepo := &repository.MockRepo{}

	mockRepo.On("UpdateCategory", mock.Anything, 1, 1, "Updated Category").
		Return(nil)

	handler := &CategoryHandler{Repo: mockRepo}

	body := map[string]string{"name": "Updated Category"}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/categories/update?id=1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))

	w := httptest.NewRecorder()
	handler.UpdateCategory(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestDeleteCategoryHandler(t *testing.T) {
	mockRepo := &repository.MockRepo{}

	mockRepo.On("DeleteCategory", mock.Anything, 1, 1).
		Return(nil)

	handler := &CategoryHandler{Repo: mockRepo}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/categories/delete?id=1", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))

	w := httptest.NewRecorder()
	handler.DeleteCategory(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}
