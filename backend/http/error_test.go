package http_test

import (
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
	"github.com/jacob-ian/jacobianmatthews.com/backend/http"
)

func TestHandleErrorCustom(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := backend.NewError(nethttp.StatusUnauthorized, "An error occurred")
	http.HandleError(recorder, err)

	res := recorder.Result()
	status := res.StatusCode

	if status != nethttp.StatusUnauthorized {
		t.Errorf("Error handler returned incorrect status code: got %v want %v", status, nethttp.StatusUnauthorized)
	}

	expected := `{"error":"unauthorized","error_description":"An error occurred"}`
	body := recorder.Body.String()
	if body != expected {
		t.Errorf("Error handler returned incorrect body: got %v want %v", body, expected)
	}
}
