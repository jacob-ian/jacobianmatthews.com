package http

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
)

type JsonDecoder struct {
	reader  io.Reader
	decoder *json.Decoder
}

// Decode a reader as JSON and throw if required fields are missing
func (d *JsonDecoder) Decode(v any) error {
	d.decoder.DisallowUnknownFields()
	err := d.decoder.Decode(v)
	if err != nil {
		return backend.NewError(http.StatusBadRequest, "Invalid JSON")
	}
	fields := reflect.ValueOf(v).Elem()
	for i := 0; i < fields.NumField(); i++ {
		tag := fields.Type().Field(i).Tag.Get("required")
		if strings.Contains(tag, "true") && fields.Field(i).IsZero() {
			return backend.NewError(http.StatusBadRequest, "Missing required field")
		}
	}
	return nil
}

// Create a JSON decoder with required field validation
func NewJsonDecoder(r io.Reader) *JsonDecoder {
	return &JsonDecoder{
		reader:  r,
		decoder: json.NewDecoder(r),
	}
}
