package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type ValidationErrorDetail struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type StructuredErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`

	Details []ValidationErrorDetail `json:"details,omitempty"`
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ERROR: Failed to marshal JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{"error":"Internal server error"}`)); err != nil {
			log.Printf("CRITICAL WRITE FAILURE: Could not write fallback error to client: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(response); err != nil {
		log.Printf("CRITICAL WARNING: Failed to write response body to client: %v", err)
	}
}

func RespondWithError(w http.ResponseWriter, status int, message string, details []ValidationErrorDetail) {
	errorPayload := StructuredErrorResponse{
		Message: message,
	}

	if len(details) > 0 {
		errorPayload.Details = details
	}

	RespondWithJSON(w, status, errorPayload)
}
