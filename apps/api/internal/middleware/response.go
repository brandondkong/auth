package middleware

import "net/http"

func WriteResponse(w http.ResponseWriter, responseCode int, message string) {
	w.WriteHeader(responseCode)
	w.Write([]byte(message))
}
