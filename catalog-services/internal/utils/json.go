package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ErrorJSON(w http.ResponseWriter, status int, message string) {
	type jsonError struct {
		Error string `json:"error"`
	}

	WriteJSON(w, status, jsonError{Error: message})
}
