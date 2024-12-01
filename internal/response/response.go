package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"massage"`
}

type SuccessResponse struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(ErrorResponse{
		StatusCode: statusCode,
		Message:    message,
	})
}

func Success(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(SuccessResponse{
		StatusCode: statusCode,
		Data:       data,
	})
}
