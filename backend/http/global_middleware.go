package http

import (
	"net/http"
	"strings"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
)

type Header struct {
	Name  string
	Value string
}

type GlobalMiddlewareConfig struct {
	CorsOrigin      string
	Accept          string
	RequiredHeaders []Header
	ResponseHeaders []Header
}

type GlobalMiddleware struct {
	handler http.Handler
	config  GlobalMiddlewareConfig
}

func (m *GlobalMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.writeResponseHeaders(w)
	e := m.checkRequestHeaders(r)
	if e != nil {
		NewResponseWriter(w, r).HandleError(e)
		return
	}
	m.handler.ServeHTTP(w, r)
}

func (m *GlobalMiddleware) writeResponseHeaders(w http.ResponseWriter) {
	responseHeaders := append(
		m.config.ResponseHeaders,
		Header{Name: "Accept", Value: m.config.Accept},
		Header{Name: "Access-Control-Allow-Origin", Value: m.config.CorsOrigin},
	)
	for _, header := range responseHeaders {
		w.Header().Set(header.Name, header.Value)
	}
}

func (m *GlobalMiddleware) checkRequestHeaders(r *http.Request) error {
	err := checkContentType(r, m.config.Accept)
	if err != nil {
		return err
	}

	for _, required := range m.config.RequiredHeaders {
		if r.Header.Get(required.Name) == required.Value {
			break
		}
		return backend.NewError(backend.BadRequestError, "Invalid Request")
	}

	return nil

}

func checkContentType(r *http.Request, accepted string) error {
	if r.Method == "GET" {
		return nil
	}
	contentType := r.Header.Get("Content-Type")
	values := strings.Split(accepted, ", ")
	matches := 0
	for _, value := range values {
		if contentType == value {
			matches += 1
		}
	}
	if matches == 1 {
		return nil
	}
	return backend.NewError(http.StatusNotAcceptable, "Not Acceptable")
}

// Create middleware that defines required request headers and the global response headers
func NewGlobalMiddleware(handler http.Handler, config GlobalMiddlewareConfig) *GlobalMiddleware {
	return &GlobalMiddleware{
		handler: handler,
		config:  config,
	}
}
