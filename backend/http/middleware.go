package http

import "net/http"

type RequestFilter struct {
	ContentType string
}

type RequestMiddleware struct {
	handler http.Handler
	filter  RequestFilter
}

func (m *RequestMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != m.filter.ContentType {
		http.Error(w, "\"Content-Type\" \"application/json\" is required", http.StatusBadRequest)
		return
	}
	m.handler.ServeHTTP(w, r)
}

// Create a request filtering middleware to wrap an HTTP handler with
func NewRequestMiddleware(handler http.Handler, filter RequestFilter) *RequestMiddleware {
	return &RequestMiddleware{
		handler: handler,
		filter:  filter,
	}
}
