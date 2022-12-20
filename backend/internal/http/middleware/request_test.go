package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/middleware"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
)

type header struct {
	Name  string
	Value string
}

type requestConfig struct {
	Headers []middleware.Header
	Method  string
	URI     string
	Body    io.Reader
}

type requestMiddlewareTestConfig struct {
	Middleware middleware.RequestMiddlewareConfig
	Request    requestConfig
}

func setupGlobalMiddlewareTest(config requestMiddlewareTestConfig) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte{0, 1})
	})
	m := middleware.NewRequestMiddleware(config.Middleware)
	m.Inject(h, res.NewResponseWriterFactory(res.ResponseWriterConfig{}))
	req := httptest.NewRequest(config.Request.Method, config.Request.URI, config.Request.Body)
	for _, header := range config.Request.Headers {
		req.Header.Set(header.Name, header.Value)
	}
	m.ServeHTTP(rr, req)
	return rr
}

func TestWriteResponseHeaders(t *testing.T) {
	rr := setupGlobalMiddlewareTest(requestMiddlewareTestConfig{
		Middleware: middleware.RequestMiddlewareConfig{
			CorsOrigin:      "abcd",
			Accept:          "application/json",
			ResponseHeaders: []middleware.Header{{Name: "Test-Header", Value: "test"}},
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
	rr := setupGlobalMiddlewareTest(requestMiddlewareTestConfig{
		Middleware: middleware.RequestMiddlewareConfig{
			CorsOrigin:      "abcd",
			Accept:          "application/json, application/grpc-web",
			ResponseHeaders: []middleware.Header{{Name: "Test-Header", Value: "test"}},
		},
		Request: requestConfig{
			Method:  "POST",
			URI:     "/test",
			Body:    nil,
			Headers: []middleware.Header{{Name: "Content-Type", Value: "application/json"}},
		}})
	res := rr.Result()
	statusCode := res.StatusCode
	if statusCode != http.StatusOK {
		t.Errorf("Received unexpected status: got '%v' want '%v', body %v", statusCode, http.StatusOK, rr.Body.String())
	}
}

func TestCheckAcceptedContentTypeBad(t *testing.T) {
	rr := setupGlobalMiddlewareTest(requestMiddlewareTestConfig{
		Middleware: middleware.RequestMiddlewareConfig{
			CorsOrigin:      "abcd",
			Accept:          "application/json, application/grpc-web",
			ResponseHeaders: []middleware.Header{{Name: "Test-Header", Value: "test"}},
		},
		Request: requestConfig{
			Method:  "POST",
			URI:     "/test",
			Body:    nil,
			Headers: []middleware.Header{{Name: "Content-Type", Value: "text/plain"}},
		}})
	res := rr.Result()
	statusCode := res.StatusCode
	if statusCode != http.StatusNotAcceptable {
		t.Errorf("Received unexpected status: got '%v' want '%v', body %v", statusCode, http.StatusNotAcceptable, rr.Body.String())
	}
}

func TestCheckRequestHeadersMissing(t *testing.T) {
	rr := setupGlobalMiddlewareTest(requestMiddlewareTestConfig{
		Middleware: middleware.RequestMiddlewareConfig{
			CorsOrigin:      "abcd",
			Accept:          "application/json, application/grpc-web",
			RequiredHeaders: []middleware.Header{{Name: "Test-Header", Value: "test"}},
		},
		Request: requestConfig{
			Method:  "POST",
			URI:     "/test",
			Body:    nil,
			Headers: []middleware.Header{{Name: "Content-Type", Value: "application/json"}},
		}})
	res := rr.Result()
	statusCode := res.StatusCode
	if statusCode != http.StatusBadRequest {
		t.Errorf("Received unexpected status: got '%v' want '%v', body %v", statusCode, http.StatusBadRequest, rr.Body.String())
	}
}

func TestCheckRequestHeadersPresent(t *testing.T) {
	rr := setupGlobalMiddlewareTest(requestMiddlewareTestConfig{
		Middleware: middleware.RequestMiddlewareConfig{
			CorsOrigin:      "abcd",
			Accept:          "application/json, application/grpc-web",
			RequiredHeaders: []middleware.Header{{Name: "Test-Header", Value: "test"}},
		},
		Request: requestConfig{
			Method:  "POST",
			URI:     "/test",
			Body:    nil,
			Headers: []middleware.Header{{Name: "Content-Type", Value: "application/json"}, {Name: "Test-Header", Value: "test"}},
		}})
	res := rr.Result()
	statusCode := res.StatusCode
	if statusCode != http.StatusOK {
		t.Errorf("Received unexpected status: got '%v' want '%v', body %v", statusCode, http.StatusOK, rr.Body.String())
	}
}
