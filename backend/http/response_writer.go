package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
)

type Afterware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}

type ResponseWriter struct {
	writer    http.ResponseWriter
	request   *http.Request
	afterware []Afterware
}

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// Write a JSON response back to the client
func (w *ResponseWriter) Write(statusCode int, res any) {
	for _, ware := range w.afterware {
		err := ware.ServeHTTP(w.writer, w.request)
		if err != nil {
			w.HandleError(err)
			return
		}
	}
	w.writer.WriteHeader(statusCode)
	json.NewEncoder(w.writer).Encode(res)
}

// Handles a thrown error and writes a JSON response to the client
func (w *ResponseWriter) HandleError(e error) {
	err, ok := e.(backend.Error)
	if !ok {
		w.WriteError(err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteError(err.GetMessage(), err.GetCode())
}

// Writes a JSON error to the client
func (w *ResponseWriter) WriteError(description string, statusCode int) {
	w.writer.Header().Set("Content-Type", "application/json")
	error := ErrorResponse{
		Error:            GetStatusError(statusCode),
		ErrorDescription: description,
	}
	w.Write(statusCode, error)
}

// Returns the formatted text error for a status code
func GetStatusError(code int) string {
	text := http.StatusText(code)
	lowercase := strings.ToLower(text)
	underscore := strings.ReplaceAll(lowercase, " ", "_")
	return underscore
}

func NewResponseWriter(w http.ResponseWriter, r *http.Request) *ResponseWriter {
	return &ResponseWriter{
		writer:    w,
		request:   r,
		afterware: []Afterware{NewCsrfAfterware()},
	}
}
