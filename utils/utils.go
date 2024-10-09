package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("JSON encoding failed: %v", err)
		return err
	}
	return nil
}

func WriteError(w http.ResponseWriter, status int, err error) {
	if jsonErr := WriteJSON(w, status, map[string]string{"error": err.Error()}); jsonErr != nil {
		http.Error(w, "An error occurred", http.StatusInternalServerError)
	}
}
