package middleware

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, responseCode int, message string) {
	w.WriteHeader(responseCode)
	w.Write([]byte(message))
}

func WriteJsonResponse(w http.ResponseWriter, responseCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	json.NewEncoder(w).Encode(data)

}
