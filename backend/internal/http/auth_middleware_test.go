package http_test

import (
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http"
	"github.com/jacob-ian/jacobianmatthews.com/backend/mock"
)

type authMiddlewareTest struct {
	Name                     string
	MockSessionServiceValues mock.MockSessionServiceValues
	RequestCookies           []nethttp.Cookie
	ExpectedStatusCode       int
	ExpectedContext          *core.SessionUser
	ExpectedCookies          []nethttp.Cookie
}

type authMiddlewareSuite struct {
	Tests []authMiddlewareTest
}

func runAuthMiddlewareSuite(t *testing.T, suite authMiddlewareSuite) {
	for i := range suite.Tests {
		test := suite.Tests[i]

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(nethttp.MethodGet, "/", nil)
		h := nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			req = r
			w.WriteHeader(200)
			w.Write([]byte("hello"))
		})

		a := mock.NewSessionService(test.MockSessionServiceValues)
		m := http.NewAuthMiddleware(h, a)

		for j := range test.RequestCookies {
			cookie := test.RequestCookies[j]
			req.AddCookie(&cookie)
		}

		m.ServeHTTP(rr, req)

		if want, got := test.ExpectedStatusCode, rr.Result().StatusCode; want != got {
			t.Errorf("'%v' failed. Unexpected status code, want %v got %v", test.Name, want, got)
		}

		userContext, ok := core.UserFromContext(req.Context())
		if ok {
			if want, got := *test.ExpectedContext, *userContext; want != got {
				t.Errorf("'%v' failed. Unexpected request context, want %v got %v", test.Name, want, got)
			}
		} else {
			if want, got := test.ExpectedContext, userContext; want != got {
				t.Errorf("'%v' failed. Unexpected request context, want %v got %v", test.Name, want, got)
			}
		}
	}
}

func TestAuthMiddleware(t *testing.T) {
	sessionUser := core.SessionUser{
		Role: core.Role{
			Id:        uuid.Must(uuid.NewRandom()),
			Name:      "Admin",
			CreatedAt: time.Now(),
		},
		User: core.User{
			Id:            uuid.Must(uuid.NewRandom()).String(),
			Name:          "User User",
			Email:         "user@email.com",
			EmailVerified: true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}
	runAuthMiddlewareSuite(t, authMiddlewareSuite{
		Tests: []authMiddlewareTest{
			{
				Name: "Should attach session user to context if has valid session cookie",
				MockSessionServiceValues: mock.MockSessionServiceValues{
					VerifySession: mock.MockResponse{
						Value: sessionUser,
						Error: nil,
					},
				},
				RequestCookies: []nethttp.Cookie{
					{
						Name:   "session",
						Value:  "session-cookie",
						MaxAge: 10,
					},
				},
				ExpectedStatusCode: 200,
				ExpectedContext:    &sessionUser,
				ExpectedCookies: []nethttp.Cookie{
					{
						Name:   "session",
						Value:  "session-cookie",
						MaxAge: 10,
					},
				},
			},
			{
				Name: "Should not attach a user when not signed in",
				MockSessionServiceValues: mock.MockSessionServiceValues{
					VerifySession: mock.MockResponse{
						Value: core.SessionUser{},
						Error: nil,
					},
				},
				RequestCookies:     []nethttp.Cookie{},
				ExpectedCookies:    []nethttp.Cookie{},
				ExpectedStatusCode: 200,
				ExpectedContext:    nil,
			},
			{
				Name: "Should respond with an error and remove the cookie if the session cookie is invalid",
				MockSessionServiceValues: mock.MockSessionServiceValues{
					VerifySession: mock.MockResponse{
						Value: core.SessionUser{},
						Error: core.NewError(core.BadRequestError, "Invalid session"),
					},
				},
				RequestCookies: []nethttp.Cookie{
					{
						Name:   "session",
						Value:  "session-cookie",
						MaxAge: 10,
					},
				},
				ExpectedStatusCode: 400,
				ExpectedContext:    nil,
				ExpectedCookies: []nethttp.Cookie{
					{
						Name:   "session",
						Value:  "",
						MaxAge: -1,
					},
				},
			},
		},
	})
}
