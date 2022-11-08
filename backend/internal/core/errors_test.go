package core_test

import (
	"testing"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
)

func TestError(t *testing.T) {
	error := core.NewError(core.InternalError, "This is a test")
	got := error.Error()
	want := "500: This is a test"
	if got != want {
		t.Errorf("Got '%v', want '%v'", got, want)
	}
}
