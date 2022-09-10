package http_test

import (
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
	"github.com/jacob-ian/jacobianmatthews.com/backend/http"
)

func TestHandleErrorCustom(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := backend.NewError(nethttp.StatusUnauthorized, "An error occurred")
	http.NewResponseWriter(recorder, &nethttp.Request{}).HandleError(err)

	res := recorder.Result()
	status := res.StatusCode

	if status != nethttp.StatusUnauthorized {
		t.Errorf("Error handler returned incorrect status code: got %v want %v", status, nethttp.StatusUnauthorized)
	}

	expected := `{"error":"unauthorized","error_description":"An error occurred"}`
	body := recorder.Body.String()
	if strings.Compare(body, expected) == 0 {
		t.Errorf("Error handler returned incorrect body: got %v want %v", body, expected)
	}
}

func TestHandleErrorUnknown(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := errors.New("Something happened")
	http.NewResponseWriter(recorder, &nethttp.Request{}).HandleError(err)

	res := recorder.Result()
	status := res.StatusCode

	expectedStatus := nethttp.StatusBadRequest
	if status != expectedStatus {
		t.Errorf("Error handler returned incorrect status code: got %v want %v", status, expectedStatus)
	}

	expected := `{"error":"bad_request","error_description":"Something happened"}`
	body := recorder.Body.String()
	if strings.Compare(body, expected) == 0 {
		t.Errorf("Error handler returned incorrect body: got %v want %v", body, expected)
	}
}
