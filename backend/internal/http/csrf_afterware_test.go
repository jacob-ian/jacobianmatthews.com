package http_test

import (
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http"
)

func TestAddCsrfToken(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	a := http.NewCsrfAfterware()
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
