package http

import (
	"encoding/json"
	"net/http"

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
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	error := ErrorResponse{
		Error:            http.StatusText(statusCode),
		ErrorDescription: description,
	}
	json.NewEncoder(w).Encode(error)
}
