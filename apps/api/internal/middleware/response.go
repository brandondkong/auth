package middleware

import (
	"encoding/json"
	"net/http"
)

type ResponseOptions[T any] struct {
	Code	int
	Error	*string
	Message	string
	Data	T
}

type JsonResponse[T any] struct {
	Success bool	`json:"success"`
	Error	*string	`json:"error"`
	Message	string	`json:"message"`
	Data	T		`json:"data"`
}

func WriteResponse(w http.ResponseWriter, responseCode int, message string) {
	w.WriteHeader(responseCode)
	w.Write([]byte(message))
}

func WriteJsonResponse[T interface{}](w http.ResponseWriter, options ResponseOptions[T]) {
	w.Header().Set("Content-Type", "application")
	w.WriteHeader(options.Code)
	json.NewEncoder(w).Encode(JsonResponse[T]{
		Success: options.Error == nil,
		Error: options.Error,
		Message: options.Message,
		Data: options.Data,
	})
}
