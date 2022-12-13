package middleware_test

import (
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/middleware"
)

func TestAddCsrfToken(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	a := middleware.NewCsrfMiddleware()
	err := a.ServeHTTP(rr, req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	res := rr.Result()

	var gotToken *nethttp.Cookie = nil
	for _, cookie := range res.Cookies() {
		if cookie.Name == "csrfToken" {
			gotToken = cookie
		}
	}

	if gotToken == nil {
		t.Errorf("Expected csrfToken Cookie: got %v", gotToken)
	}
}
