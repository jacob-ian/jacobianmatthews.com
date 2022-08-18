package backend_test

import (
	"testing"

	"github.com/jacob-ian/jacobianmatthews.com/backend"
)

func TestError(t *testing.T) {
	error := backend.NewError(backend.InternalError, "This is a test")
	got := error.Error()
	want := "500: This is a test"
	if got != want {
		t.Errorf("Got '%v', want '%v'", got, want)
	}
}
