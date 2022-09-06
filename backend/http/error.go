package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
)

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// Handle errors thrown in the HTTP application
func HandleError(w http.ResponseWriter, e error) {
	err, ok := e.(backend.Error)
	if !ok {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeError(w, err.GetMessage(), err.GetCode())
}

// Responds with a JSON error
func writeError(w http.ResponseWriter, description string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	error := ErrorResponse{
		Error:            GetStatusError(statusCode),
		ErrorDescription: description,
	}
	json.NewEncoder(w).Encode(error)
}

// Returns the formatted text error for a status code
func GetStatusError(code int) string {
	text := http.StatusText(code)
	lowercase := strings.ToLower(text)
	underscore := strings.ReplaceAll(lowercase, " ", "_")
	return underscore
}