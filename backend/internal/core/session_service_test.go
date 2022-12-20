package core_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/mock"
)

type sessionServiceTest struct {
	Name                  string
	AuthProviderValues    mock.MockAuthProviderValues
	UserRespositoryValues mock.MockUserRepositoryValues
	AuthServiceValues     mock.MockAuthServiceValues
	ExpectedOutput        any
	ExpectedError         error
}

type sessionServiceSuite struct {
	Func  func(service *core.CoreSessionService) (any, error)
	Tests []sessionServiceTest
}

func runSessionServiceSuite(t *testing.T, suite sessionServiceSuite) {
	tests := suite.Tests
	for i := range tests {
		test := tests[i]

		service := core.NewSessionService(core.CoreSessionServiceConfig{
			AuthProvider:   mock.NewAuthProvider(test.AuthProviderValues),
			UserRepository: mock.NewUserRepository(test.UserRespositoryValues),
			AuthService:    mock.NewAuthService(test.AuthServiceValues),
		})

		res, err := suite.Func(service)

		if want, got := test.ExpectedError, err; want != got {
			t.Errorf("'%v' failed. Unexpected error, want %v got %v", test.Name, want, got)
		}

		if want, got := test.ExpectedOutput, res; want != got {
			t.Errorf("'%v' failed. Unexpected output, want %v got %v", test.Name, want, got)
		}
	}
}

func TestSessionService_StartSession(t *testing.T) {
	runSessionServiceSuite(t, sessionServiceSuite{
		Func: func(s *core.CoreSessionService) (any, error) {
			return s.StartSession(context.Background(), "idToken")
		},
		Tests: []sessionServiceTest{
			{
				Name: "Should all pass",
				AuthProviderValues: mock.MockAuthProviderValues{
					VerifyIdToken: mock.MockResponse{
						Value: &core.Token{
							Subject: "user1",
							Claims: map[string]any{
								"auth_time": time.Now().Unix() - int64(time.Minute),
							}},
						Error: nil,
					},
					CreateSessionCookie: mock.MockResponse{
						Value: "session-cookie",
						Error: nil,
					},
					GetUserDetails: mock.MockResponse{
						Value: core.AuthProviderUser{},
						Error: nil,
					},
				},
				UserRespositoryValues: mock.MockUserRepositoryValues{
					FindById: mock.MockResponse{
						Value: core.User{},
						Error: nil,
					},
					Create: mock.MockResponse{
						Value: core.User{},
						Error: nil,
					},
				},
				AuthServiceValues: mock.MockAuthServiceValues{
					GiveUserRoleByName: mock.MockResponse{
						Value: core.Role{},
						Error: nil,
					},
				},
				ExpectedError: nil,
				ExpectedOutput: core.Session{
					Cookie:    "session-cookie",
					ExpiresIn: time.Minute * 15,
				},
			},
			{
				Name: "Should fail if request has invalid ID Token",
				AuthProviderValues: mock.MockAuthProviderValues{
					VerifyIdToken: mock.MockResponse{
						Value: &core.Token{},
						Error: errors.New("This is a dodgy token!"),
					},
					CreateSessionCookie: mock.MockResponse{
						Value: "session-cookie",
						Error: nil,
					},
					GetUserDetails: mock.MockResponse{
						Value: core.AuthProviderUser{},
						Error: nil,
					},
				},
				UserRespositoryValues: mock.MockUserRepositoryValues{
					FindById: mock.MockResponse{
						Value: core.User{},
						Error: nil,
					},
					Create: mock.MockResponse{
						Value: core.User{},
						Error: nil,
					},
				},
				AuthServiceValues: mock.MockAuthServiceValues{
					GiveUserRoleByName: mock.MockResponse{
						Value: core.Role{},
						Error: nil,
					},
				},
				ExpectedError:  core.NewError(core.BadRequestError, "Invalid ID Token"),
				ExpectedOutput: core.Session{},
			},
			{
				Name: "Should fail if user authenticated more than 5 minutes ago",
				AuthProviderValues: mock.MockAuthProviderValues{
					VerifyIdToken: mock.MockResponse{
						Value: &core.Token{
							Subject: "user1",
							Claims: map[string]any{
								"auth_time": time.Now().Unix() - int64(time.Minute*6),
							}},
						Error: nil,
					},
					CreateSessionCookie: mock.MockResponse{
						Value: "session-cookie",
						Error: nil,
					},
					GetUserDetails: mock.MockResponse{
						Value: core.AuthProviderUser{},
						Error: nil,
					},
				},
				UserRespositoryValues: mock.MockUserRepositoryValues{
					FindById: mock.MockResponse{
						Value: core.User{},
						Error: nil,
					},
					Create: mock.MockResponse{
						Value: core.User{},
						Error: nil,
					},
				},
				AuthServiceValues: mock.MockAuthServiceValues{
					GiveUserRoleByName: mock.MockResponse{
						Value: core.Role{},
						Error: nil,
					},
				},
				ExpectedOutput: core.Session{},
				ExpectedError:  core.NewError(core.UnauthenticatedError, "Unauthenticated"),
			},
			{
				Name: "Should fail if user repository throws a non-NotFound error",
				AuthProviderValues: mock.MockAuthProviderValues{
					VerifyIdToken: mock.MockResponse{
						Value: &core.Token{
							Subject: "user1",
							Claims: map[string]any{
								"auth_time": time.Now().Unix() - int64(time.Minute),
							}},
						Error: nil,
					},
					CreateSessionCookie: mock.MockResponse{
						Value: "session-cookie",
						Error: nil,
					},
					GetUserDetails: mock.MockResponse{
						Value: core.AuthProviderUser{},
						Error: nil,
					},
				},
				UserRespositoryValues: mock.MockUserRepositoryValues{
					FindById: mock.MockResponse{
						Value: core.User{},
						Error: core.NewError(core.InternalError, "Something went wrong"),
					},
					Create: mock.MockResponse{
						Value: core.User{},
						Error: nil,
					},
				},
				AuthServiceValues: mock.MockAuthServiceValues{
					GiveUserRoleByName: mock.MockResponse{
						Value: core.Role{},
						Error: nil,
					},
				},
				ExpectedOutput: core.Session{},
				ExpectedError:  core.NewError(core.InternalError, "Something went wrong"),
			},
			{
				Name: "Should fail if user repository throws a NotFound error and provider user details throws an error",
				AuthProviderValues: mock.MockAuthProviderValues{
					VerifyIdToken: mock.MockResponse{
						Value: &core.Token{
							Subject: "user1",
							Claims: map[string]any{
								"auth_time": time.Now().Unix() - int64(time.Minute),
							}},
						Error: nil,
					},
					CreateSessionCookie: mock.MockResponse{
						Value: "session-cookie",
						Error: nil,
					},
					GetUserDetails: mock.MockResponse{
						Value: core.AuthProviderUser{},
						Error: errors.New("Couldn't find user with that ID"),
					},
				},
				UserRespositoryValues: mock.MockUserRepositoryValues{
					FindById: mock.MockResponse{
						Value: core.User{},
						Error: core.NewError(core.NotFoundError, "Couldn't find user"),
					},
					Create: mock.MockResponse{
						Value: core.User{},
						Error: nil,
					},
				},
				AuthServiceValues: mock.MockAuthServiceValues{
					GiveUserRoleByName: mock.MockResponse{
						Value: core.Role{},
						Error: nil,
					},
				},
				ExpectedOutput: core.Session{},
				ExpectedError:  core.NewError(core.BadRequestError, "Invalid Firebase User ID"),
			},
			{
				Name: "Should fail if user repository throws an error on Create",
				AuthProviderValues: mock.MockAuthProviderValues{
					VerifyIdToken: mock.MockResponse{
						Value: &core.Token{
							Subject: "user1",
							Claims: map[string]any{
								"auth_time": time.Now().Unix() - int64(time.Minute),
							}},
						Error: nil,
					},
					CreateSessionCookie: mock.MockResponse{
						Value: "session-cookie",
						Error: nil,
					},
					GetUserDetails: mock.MockResponse{
						Value: core.AuthProviderUser{
							Id:            uuid.Must(uuid.NewRandom()).String(),
							DisplayName:   "Namey Name",
							Email:         "email@email.com",
							EmailVerified: false,
							ImageUrl:      "",
						},
						Error: nil,
					},
				},
				UserRespositoryValues: mock.MockUserRepositoryValues{
					FindById: mock.MockResponse{
						Value: core.User{},
						Error: core.NewError(core.NotFoundError, "Couldn't find user"),
					},
					Create: mock.MockResponse{
						Value: core.User{},
						Error: core.NewError(core.InternalError, "Couldn't create user"),
					},
				},
				AuthServiceValues: mock.MockAuthServiceValues{
					GiveUserRoleByName: mock.MockResponse{
						Value: core.Role{},
						Error: nil,
					},
				},
				ExpectedOutput: core.Session{},
				ExpectedError:  core.NewError(core.InternalError, "Couldn't create user"),
			},
			{
				Name: "Should fail if auth service can't give new user a role",
				AuthProviderValues: mock.MockAuthProviderValues{
					VerifyIdToken: mock.MockResponse{
						Value: &core.Token{
							Subject: "user1",
							Claims: map[string]any{
								"auth_time": time.Now().Unix() - int64(time.Minute),
							}},
						Error: nil,
					},
					CreateSessionCookie: mock.MockResponse{
						Value: "session-cookie",
						Error: nil,
					},
					GetUserDetails: mock.MockResponse{
						Value: core.AuthProviderUser{
							Id:            uuid.Must(uuid.NewRandom()).String(),
							DisplayName:   "Namey Name",
							Email:         "email@email.com",
							EmailVerified: false,
							ImageUrl:      "",
						},
						Error: nil,
					},
				},
				UserRespositoryValues: mock.MockUserRepositoryValues{
					FindById: mock.MockResponse{
						Value: core.User{},
						Error: core.NewError(core.NotFoundError, "Couldn't find user"),
					},
					Create: mock.MockResponse{
						Value: core.User{},
						Error: nil,
					},
				},
				AuthServiceValues: mock.MockAuthServiceValues{
					GiveUserRoleByName: mock.MockResponse{
						Value: core.Role{},
						Error: core.NewError(core.InternalError, "Couldn't give user role"),
					},
				},
				ExpectedOutput: core.Session{},
				ExpectedError:  core.NewError(core.InternalError, "Couldn't give user role"),
			},
			{
				Name: "Should fail if provider fails to create a session cookie",
				AuthProviderValues: mock.MockAuthProviderValues{
					VerifyIdToken: mock.MockResponse{
						Value: &core.Token{
							Subject: "user1",
							Claims: map[string]any{
								"auth_time": time.Now().Unix() - int64(time.Minute),
							}},
						Error: nil,
					},
					CreateSessionCookie: mock.MockResponse{
						Value: "",
						Error: errors.New("Couldn't create session cookie"),
					},
					GetUserDetails: mock.MockResponse{
						Value: core.AuthProviderUser{
							Id:            uuid.Must(uuid.NewRandom()).String(),
							DisplayName:   "Namey Name",
							Email:         "email@email.com",
							EmailVerified: false,
							ImageUrl:      "",
						},
						Error: nil,
					},
				},
				UserRespositoryValues: mock.MockUserRepositoryValues{
					FindById: mock.MockResponse{
						Value: core.User{},
						Error: core.NewError(core.NotFoundError, "Couldn't find user"),
					},
					Create: mock.MockResponse{
						Value: core.User{},
						Error: nil,
					},
				},
				AuthServiceValues: mock.MockAuthServiceValues{
					GiveUserRoleByName: mock.MockResponse{
						Value: core.Role{},
						Error: nil,
					},
				},
				ExpectedOutput: core.Session{},
				ExpectedError:  core.NewError(core.InternalError, "An error occurred whilst signing in"),
			},
		},
	})
}

func TestSessionService_VerifySession(t *testing.T) {
	role := core.Role{
		Id:        uuid.Must(uuid.NewRandom()),
		Name:      "Admin",
		CreatedAt: time.Now().Add(-time.Hour * 120),
	}
	user := core.User{
		Id:            uuid.Must(uuid.NewRandom()).String(),
		Name:          "Usery User",
		Email:         "emaily@email.com",
		EmailVerified: true,
		CreatedAt:     time.Now().Add(-time.Hour * 72),
		UpdatedAt:     time.Now().Add(-time.Hour * 72),
	}

	runSessionServiceSuite(t, sessionServiceSuite{
		Func: func(service *core.CoreSessionService) (any, error) {
			return service.VerifySession(context.Background(), "session-cookie")
		},
		Tests: []sessionServiceTest{
			{
				Name: "Should return the user's role and user details",
				AuthProviderValues: mock.MockAuthProviderValues{
					VerifySessionCookie: mock.MockResponse{
						Value: &core.Token{
							Subject: user.Id,
						},
						Error: nil,
					},
				},
				AuthServiceValues: mock.MockAuthServiceValues{
					GetOrGiveUserRole: mock.MockResponse{
						Value: role,
						Error: nil,
					},
				},
				UserRespositoryValues: mock.MockUserRepositoryValues{
					FindById: mock.MockResponse{
						Value: user,
						Error: nil,
					},
				},
				ExpectedOutput: core.SessionUser{
					User: user,
					Role: role,
				},
				ExpectedError: nil,
			},
			{
				Name: "Should throw an error if the session cookie is invalid",
				AuthProviderValues: mock.MockAuthProviderValues{
					VerifySessionCookie: mock.MockResponse{
						Value: &core.Token{},
						Error: errors.New("Invalid sessoin cookie dude"),
					},
				},
				AuthServiceValues: mock.MockAuthServiceValues{
					GetOrGiveUserRole: mock.MockResponse{
						Value: role,
						Error: nil,
					},
				},
				UserRespositoryValues: mock.MockUserRepositoryValues{
					FindById: mock.MockResponse{
						Value: user,
						Error: nil,
					},
				},
				ExpectedOutput: core.SessionUser{},
				ExpectedError:  core.NewError(core.UnauthenticatedError, "Invalid session"),
			},
			{
				Name: "Should throw an error if the user respository lookup fails",
				AuthProviderValues: mock.MockAuthProviderValues{
					VerifySessionCookie: mock.MockResponse{
						Value: &core.Token{
							Subject: user.Id,
						},
						Error: nil,
					},
				},
				AuthServiceValues: mock.MockAuthServiceValues{
					GetOrGiveUserRole: mock.MockResponse{
						Value: role,
						Error: nil,
					},
				},
				UserRespositoryValues: mock.MockUserRepositoryValues{
					FindById: mock.MockResponse{
						Value: core.User{},
						Error: core.NewError(core.NotFoundError, "Could not find user"),
					},
				},
				ExpectedOutput: core.SessionUser{},
				ExpectedError:  core.NewError(core.InternalError, "Could not get signed in user"),
			},
			{
				Name: "Should throw an error if the auth service role lookup fails",
				AuthProviderValues: mock.MockAuthProviderValues{
					VerifySessionCookie: mock.MockResponse{
						Value: &core.Token{
							Subject: user.Id,
						},
						Error: nil,
					},
				},
				AuthServiceValues: mock.MockAuthServiceValues{
					GetOrGiveUserRole: mock.MockResponse{
						Value: core.Role{},
						Error: core.NewError(core.NotFoundError, "Could not find user role"),
					},
				},
				UserRespositoryValues: mock.MockUserRepositoryValues{
					FindById: mock.MockResponse{
						Value: user,
						Error: nil,
					},
				},
				ExpectedOutput: core.SessionUser{},
				ExpectedError:  core.NewError(core.InternalError, "Could not verify session"),
			},
		},
	})

}
