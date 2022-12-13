package res

import "net/http"

type Afterware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}

type ResponseWriterFactory struct {
	afterware []Afterware
}

type ResponseWriterConfig struct {
	Afterware []Afterware
}

// Creates a response writer with a set of middlewares to be ran before sending a response
func NewResponseWriterFactory(config ResponseWriterConfig) *ResponseWriterFactory {
	return &ResponseWriterFactory{
		afterware: config.Afterware,
	}
}

// Create a response writer given a http request writer and request
func (f *ResponseWriterFactory) NewResponseWriter(w http.ResponseWriter, r *http.Request) *ResponseWriter {
	return &ResponseWriter{
		writer:    w,
		request:   r,
		afterware: f.afterware,
	}
}
