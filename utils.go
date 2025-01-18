package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)

	enc := json.NewEncoder(w)

	if err := enc.Encode(v); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
	}
}

func throwError(w http.ResponseWriter, message string, status int) {

	writeJSON(w, status, map[string]string{"error": message})
}
