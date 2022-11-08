package http_test

import (
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http"
	"github.com/jacob-ian/jacobianmatthews.com/backend/mock"
)

type amTestConfig struct {
	SessionCookie       string
	VerifySessionOutput mock.MockVerifySessionOutput
}

func setupAuthMiddlewareTest(config amTestConfig) (*httptest.ResponseRecorder, *nethttp.Request) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodGet, "/", nil)
	h := nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		req = r
		w.WriteHeader(200)
		w.Write([]byte("hello"))
	})
	a := mock.NewAuthService(mock.AuthServiceOutput{
		VerifySession: config.VerifySessionOutput,
	})
	m := http.NewAuthMiddleware(h, a)
	if config.SessionCookie != "" {
		req.AddCookie(&nethttp.Cookie{
			Name:   "session",
			Value:  config.SessionCookie,
			MaxAge: 10,
		})
	}
	m.ServeHTTP(rr, req)
	return rr, req
}

// When user is not signed in, it should not attach a user to context
func TestNotSignedInPassThrough(t *testing.T) {
	rr, req := setupAuthMiddlewareTest(amTestConfig{
		SessionCookie: "",
	})

	userCtx, ok := req.Context().Value(core.UserContextKey).(*core.SessionUser)
	if ok == true {
		t.Errorf("Unexpected user in context: got %v want %v", userCtx, nil)
	}

	status := rr.Result().StatusCode
	if status != nethttp.StatusOK {
		t.Errorf("Unexpected status code: got %v want %v", status, nethttp.StatusOK)
	}
}

// Should return an error code and request should not have user context
func TestSessionVerifyErrorBlock(t *testing.T) {
	rr, req := setupAuthMiddlewareTest(amTestConfig{
		SessionCookie: "coookie",
		VerifySessionOutput: mock.MockVerifySessionOutput{
			Value: nil,
			Error: core.NewError(core.BadRequestError, "Invalid session"),
		},
	})

	status := rr.Result().StatusCode
	if status != core.BadRequestError {
		t.Errorf("Unexpected status code: got %v want %v", status, core.BadRequestError)
	}

	userCtx, ok := req.Context().Value(core.UserContextKey).(*core.SessionUser)
	if ok == true {
		t.Errorf("Unexpected user in context: got %v want %v", userCtx, nil)
	}
}

// Should attach user to request context and pass through to handler
func TestValidSessionUserContext(t *testing.T) {
	user := &core.SessionUser{
		Admin: true,
		User: core.User{
			Id:        "id",
			Name:      "lolname",
			Email:     "lol@email",
			ImageUrl:  "img",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}

	rr, req := setupAuthMiddlewareTest(amTestConfig{
		SessionCookie: "coookie",
		VerifySessionOutput: mock.MockVerifySessionOutput{
			Value: user,
			Error: nil,
		},
	})

	userCtx, ok := req.Context().Value(core.UserContextKey).(*core.SessionUser)
	if ok == false || userCtx != user {
		t.Errorf("Unexpected user in context: got %v want %v", userCtx, user)
	}

	status := rr.Result().StatusCode
	if status != nethttp.StatusOK {
		t.Errorf("Unexpected status code: got %v want %v", status, nethttp.StatusOK)
	}
}
