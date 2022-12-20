package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/middleware"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/http/res"
	"github.com/jacob-ian/jacobianmatthews.com/backend/mock"
)

type sessionMiddlewareTest struct {
	Name                     string
	MockSessionServiceValues mock.MockSessionServiceValues
	RequestCookies           []http.Cookie
	ExpectedStatusCode       int
	ExpectedContext          *core.SessionUser
	ExpectedCookies          []http.Cookie
	ExpectedSetCookies       []http.Cookie
}

type sessionMiddlewareSuite struct {
	Tests []sessionMiddlewareTest
}

func runSessionMiddlewareSuite(t *testing.T, suite sessionMiddlewareSuite) {
	for i := range suite.Tests {
		test := suite.Tests[i]

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req = r
			w.WriteHeader(200)
			w.Write([]byte("hello"))
		})

		a := mock.NewSessionService(test.MockSessionServiceValues)
		m := middleware.NewSessionMiddleware(middleware.SessionMiddlewareConfig{
			SessionService: a,
		})

		m.Inject(h, res.NewResponseWriterFactory(res.ResponseWriterConfig{}))

		for j := range test.RequestCookies {
			cookie := test.RequestCookies[j]
			req.AddCookie(&cookie)
		}

		m.ServeHTTP(rr, req)

		result := rr.Result()

		if want, got := test.ExpectedStatusCode, result.StatusCode; want != got {
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

		setCookies := result.Cookies()
		for k := range test.ExpectedSetCookies {
			expected := test.ExpectedSetCookies[k]
			var actual http.Cookie
			for j := range setCookies {
				if cookie := setCookies[j]; cookie.Name == expected.Name {
					actual = *cookie
				}
			}

			if want, got := expected.Name, actual.Name; want != got {
				t.Errorf("'%v' failed. Expected cookie %v", test.Name, want)
			}
			if want, got := expected.Value, actual.Value; want != got {
				t.Errorf("'%v' failed. Unexpected cookie value, want %v got %v", test.Name, want, got)
			}
			if want, got := expected.MaxAge, actual.MaxAge; want != got {
				t.Errorf("'%v' failed. Unexpected cookie MaxAge, want %v got %v", test.Name, want, got)
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
	runSessionMiddlewareSuite(t, sessionMiddlewareSuite{
		Tests: []sessionMiddlewareTest{
			{
				Name: "Should attach session user to context and extend session cookie 15 mins if has valid session cookie",
				MockSessionServiceValues: mock.MockSessionServiceValues{
					VerifySession: mock.MockResponse{
						Value: sessionUser,
						Error: nil,
					},
				},
				RequestCookies: []http.Cookie{
					{
						Name:   "session",
						Value:  "session-cookie",
						MaxAge: 10,
					},
				},
				ExpectedStatusCode: 200,
				ExpectedContext:    &sessionUser,
				ExpectedCookies: []http.Cookie{
					{
						Name:   "session",
						Value:  "session-cookie",
						MaxAge: 10,
					},
				},
				ExpectedSetCookies: []http.Cookie{
					{
						Name:   "session",
						Value:  "session-cookie",
						MaxAge: 60 * 15,
					},
				},
			},
			{
				Name: "Should not attach a user to context when not signed in",
				MockSessionServiceValues: mock.MockSessionServiceValues{
					VerifySession: mock.MockResponse{
						Value: core.SessionUser{},
						Error: nil,
					},
				},
				RequestCookies:     []http.Cookie{},
				ExpectedCookies:    []http.Cookie{},
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
				RequestCookies: []http.Cookie{
					{
						Name:   "session",
						Value:  "session-cookie",
						MaxAge: 10,
					},
				},
				ExpectedStatusCode: 400,
				ExpectedContext:    nil,
				ExpectedCookies: []http.Cookie{
					{
						Name:   "session",
						Value:  "session-cookie",
						MaxAge: 10,
					},
				},
				ExpectedSetCookies: []http.Cookie{
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
