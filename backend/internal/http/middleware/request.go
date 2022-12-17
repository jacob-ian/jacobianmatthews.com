package middleware

import (
	"net/http"
	"strings"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
)

type Header struct {
	Name  string
	Value string
}

type RequestMiddlewareConfig struct {
	CorsOrigin      string
	Accept          string
	RequiredHeaders []Header
	ResponseHeaders []Header
}

type RequestMiddleware struct {
	handler http.Handler
	res     *res.ResponseWriterFactory
	config  RequestMiddlewareConfig
}

func (m *RequestMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.writeResponseHeaders(w)
	if e := m.checkRequestHeaders(r); e != nil {
		m.res.NewResponseWriter(w, r).HandleError(e)
		return
	}
	m.handler.ServeHTTP(w, r)
}

func (m *RequestMiddleware) writeResponseHeaders(w http.ResponseWriter) {
	responseHeaders := append(
		m.config.ResponseHeaders,
		Header{Name: "Accept", Value: m.config.Accept},
		Header{Name: "Access-Control-Allow-Origin", Value: m.config.CorsOrigin},
	)
	for _, header := range responseHeaders {
		w.Header().Set(header.Name, header.Value)
	}
}

func (m *RequestMiddleware) checkRequestHeaders(r *http.Request) error {
	err := checkContentType(r, m.config.Accept)
	if err != nil {
		return err
	}

	for _, required := range m.config.RequiredHeaders {
		if r.Header.Get(required.Name) == required.Value {
			break
		}
		return core.NewError(core.BadRequestError, "Invalid Request")
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
	return core.NewError(http.StatusNotAcceptable, "Not Acceptable")
}

func (rm *RequestMiddleware) Inject(handler http.Handler, writer *res.ResponseWriterFactory) http.Handler {
	rm.handler = handler
	rm.res = writer
	return rm
}

// Create middleware that defines required request headers and the global response headers
func NewRequestMiddleware(config RequestMiddlewareConfig) *RequestMiddleware {
	return &RequestMiddleware{
		config: config,
	}
}
