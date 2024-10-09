package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}

func WriteError(w http.ResponseWriter, status int, err error) {
	log.Printf("Error: %v", err)
	if jsonErr := WriteJSON(w, status, map[string]string{"error": err.Error()}); jsonErr != nil {
		log.Printf("JSON writing error: %v", jsonErr)
		http.Error(w, "An error occurred", http.StatusInternalServerError)
	}
}
