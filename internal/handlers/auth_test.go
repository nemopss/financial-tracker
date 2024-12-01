package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nemopss/financial-tracker/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterHandler(t *testing.T) {
	mockRepo := &repository.MockRepo{}

	mockRepo.On("CreateUser", mock.Anything, "testuser", mock.Anything).
		Return(1, nil)

	handler := &AuthHandler{
		Repo: mockRepo,
	}

	body := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Register(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var wrappedResponse struct {
		StatusCode int
		Data       map[string]interface{}
	}
	err := json.NewDecoder(resp.Body).Decode(&wrappedResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Equal(t, float64(1), wrappedResponse.Data["id"].(float64))
	assert.Equal(t, "testuser", wrappedResponse.Data["username"])

	mockRepo.AssertExpectations(t)
}

func TestRegisterHandlerInvalidPayload(t *testing.T) {
	handler := &AuthHandler{}

	body := map[string]string{
		"username": "",
		"password": "",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Register(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLoginHandler(t *testing.T) {
	mockRepo := &repository.MockRepo{}

	// Хэш для пароля "password123"
	mockUser := &repository.User{
		ID:       1,
		Username: "testuser",
		Password: "$2a$10$YSAEmDcghBgOyabYSG0cue53CAa/G1ZFP64UoayBp.Xz9rGusuPOC",
	}

	mockRepo.On("GetUserByUsername", mock.Anything, "testuser").
		Return(mockUser, nil)

	handler := &AuthHandler{
		Repo:      mockRepo,
		JWTSecret: "secret",
	}

	body := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Login(w, req)

	resp := w.Result()

	// Проверяем статус-код
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var wrappedResponse struct {
		StatusCode int
		Data       map[string]string
	}
	err := json.NewDecoder(resp.Body).Decode(&wrappedResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.NotEmpty(t, wrappedResponse.Data["token"])

	mockRepo.AssertExpectations(t)
}

func TestLoginHandlerInvalidPassword(t *testing.T) {
	mockRepo := &repository.MockRepo{}

	mockUser := &repository.User{
		ID:       1,
		Username: "testuser",
		Password: "$2a$10$eImiTXuWVxfM37uY4JANjQe7/J1mDTG6gseBPZET6n1gqKPv4rP9m", // Хэш для "password123"
	}

	mockRepo.On("GetUserByUsername", mock.Anything, "testuser").
		Return(mockUser, nil)

	handler := &AuthHandler{
		Repo:      mockRepo,
		JWTSecret: "secret",
	}

	body := map[string]string{
		"username": "testuser",
		"password": "wrongpassword",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Login(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}
