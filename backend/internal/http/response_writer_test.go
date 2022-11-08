package http_test

import (
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http"
)

func TestWriteJson(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	type response struct {
		Message string    `json:"messsage"`
		Time    time.Time `json:"time"`
	}

	now := time.Now().UTC()

	http.NewResponseWriter(rr, req).Write(nethttp.StatusCreated, response{
		Message: "Hello",
		Time:    now,
	})

	res := rr.Result()
	status := res.StatusCode

	if status != nethttp.StatusCreated {
		t.Errorf("Unexpected status: got %v want %v", status, nethttp.StatusCreated)
	}

	nowFmt, err := json.Marshal(now)
	if err != nil {
		t.Errorf("Unexpected json error: %v", err)
	}

	got := rr.Body.String()
	want := `{"message": "Hello","time":` + string(nowFmt) + `}`

	if strings.Compare(got, want) == 0 {
		t.Errorf("Unexpected body: got %v want %v", got, want)
	}
}

func TestWriteError(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	http.NewResponseWriter(rr, req).WriteError("Test Error", nethttp.StatusMethodNotAllowed)

	statusCode := rr.Result().StatusCode
	if statusCode != nethttp.StatusMethodNotAllowed {
		t.Errorf("Unexpected status: got %v, want %v", statusCode, nethttp.StatusMethodNotAllowed)
	}

	contentType := rr.Result().Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Unexpected content type: got %v want %v", contentType, "application/json")
	}

	got := rr.Body.String()
	want := `{"error":"method_not_allowed","error_description":"Test Error"}`
	if strings.Compare(got, want) == 0 {
		t.Errorf("Unexpected body: got %v want %v", got, want)
	}
}

func TestHandleErrorCustomType(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	err := core.NewError(core.InternalError, "A test error")

	http.NewResponseWriter(rr, req).HandleError(err)

	sc := rr.Result().StatusCode
	if sc != core.InternalError {
		t.Errorf("Unexpected status: got %v want %v", sc, core.InternalError)
	}

	ct := rr.Result().Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Unexpected content type: got %v want %v", ct, "application/json")
	}

	got := rr.Body.String()
	want := `{"error":"internal_server_error","error_description":"A test error"}`
	if strings.Compare(got, want) == 0 {
		t.Errorf("Unexpected body: got %v want %v", got, want)
	}
}

func TestHandleErrorUnknownType(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	err := errors.New("Something went wrong")

	http.NewResponseWriter(rr, req).HandleError(err)

	sc := rr.Result().StatusCode
	if sc != nethttp.StatusInternalServerError {
		t.Errorf("Unexpected status: got %v want %v", sc, nethttp.StatusInternalServerError)
	}

	ct := rr.Result().Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Unexpected content type: got %v want %v", ct, "application/json")
	}

	got := rr.Body.String()
	want := `{"error":"internal_server_error","error_description":"Something went wrong"}`
	if strings.Compare(got, want) == 0 {
		t.Errorf("Unexpected body: got %v want %v", got, want)
	}
}
