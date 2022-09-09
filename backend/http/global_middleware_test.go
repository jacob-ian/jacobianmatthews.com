package http_test

import (
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/jacob-ian/jacobianmatthews.com/backend/http"
)

type requestConfig struct {
	Headers []http.Header
	Method  string
	URI     string
	Body    io.Reader
}

type testConfig struct {
	Middleware http.GlobalMiddlewareConfig
	Request    requestConfig
}

func setupTest(config testConfig) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	h := nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(200)
		w.Write([]byte{0, 1})
	})
	m := http.NewGlobalMiddleware(h, config.Middleware)
	req := httptest.NewRequest(config.Request.Method, config.Request.URI, config.Request.Body)
	for _, header := range config.Request.Headers {
		req.Header.Set(header.Name, header.Value)
	}
	m.ServeHTTP(rr, req)
	return rr
}

func TestWriteResponseHeaders(t *testing.T) {
	rr := setupTest(testConfig{
		Middleware: http.GlobalMiddlewareConfig{
			CorsOrigin:      "abcd",
			Accept:          "application/json",
			ResponseHeaders: []http.Header{{Name: "Test-Header", Value: "test"}},
		},
		Request: requestConfig{
			Method: "GET",
			URI:    "/test",
			Body:   nil,
		}})

	res := rr.Result()

	actual := []string{
		res.Header.Get("Access-Control-Allow-Origin"),
		res.Header.Get("Test-Header"),
	}

	expected := []string{
		"abcd",
		"test",
	}

	for i, header := range actual {
		if header != expected[i] {
			t.Errorf("Received unexpected header: got '%v' want '%v'", header, expected[i])
		}
	}
}

func TestCheckAcceptedContentTypeCorrect(t *testing.T) {
	rr := setupTest(testConfig{
		Middleware: http.GlobalMiddlewareConfig{
			CorsOrigin:      "abcd",
			Accept:          "application/json, application/grpc-web",
			ResponseHeaders: []http.Header{{Name: "Test-Header", Value: "test"}},
		},
		Request: requestConfig{
			Method:  "POST",
			URI:     "/test",
			Body:    nil,
			Headers: []http.Header{{Name: "Content-Type", Value: "application/json"}},
		}})
	res := rr.Result()
	statusCode := res.StatusCode
	if statusCode != nethttp.StatusOK {
		t.Errorf("Received unexpected status: got '%v' want '%v', body %v", statusCode, nethttp.StatusOK, rr.Body.String())
	}
}

func TestCheckAcceptedContentTypeBad(t *testing.T) {
	rr := setupTest(testConfig{
		Middleware: http.GlobalMiddlewareConfig{
			CorsOrigin:      "abcd",
			Accept:          "application/json, application/grpc-web",
			ResponseHeaders: []http.Header{{Name: "Test-Header", Value: "test"}},
		},
		Request: requestConfig{
			Method:  "POST",
			URI:     "/test",
			Body:    nil,
			Headers: []http.Header{{Name: "Content-Type", Value: "text/plain"}},
		}})
	res := rr.Result()
	statusCode := res.StatusCode
	if statusCode != nethttp.StatusNotAcceptable {
		t.Errorf("Received unexpected status: got '%v' want '%v', body %v", statusCode, nethttp.StatusNotAcceptable, rr.Body.String())
	}
}

func TestCheckRequestHeadersMissing(t *testing.T) {
	rr := setupTest(testConfig{
		Middleware: http.GlobalMiddlewareConfig{
			CorsOrigin:      "abcd",
			Accept:          "application/json, application/grpc-web",
			RequiredHeaders: []http.Header{{Name: "Test-Header", Value: "test"}},
		},
		Request: requestConfig{
			Method:  "POST",
			URI:     "/test",
			Body:    nil,
			Headers: []http.Header{{Name: "Content-Type", Value: "application/json"}},
		}})
	res := rr.Result()
	statusCode := res.StatusCode
	if statusCode != nethttp.StatusBadRequest {
		t.Errorf("Received unexpected status: got '%v' want '%v', body %v", statusCode, nethttp.StatusBadRequest, rr.Body.String())
	}
}

func TestCheckRequestHeadersPresent(t *testing.T) {
	rr := setupTest(testConfig{
		Middleware: http.GlobalMiddlewareConfig{
			CorsOrigin:      "abcd",
			Accept:          "application/json, application/grpc-web",
			RequiredHeaders: []http.Header{{Name: "Test-Header", Value: "test"}},
		},
		Request: requestConfig{
			Method:  "POST",
			URI:     "/test",
			Body:    nil,
			Headers: []http.Header{{Name: "Content-Type", Value: "application/json"}, {Name: "Test-Header", Value: "test"}},
		}})
	res := rr.Result()
	statusCode := res.StatusCode
	if statusCode != nethttp.StatusOK {
		t.Errorf("Received unexpected status: got '%v' want '%v', body %v", statusCode, nethttp.StatusOK, rr.Body.String())
	}
}
