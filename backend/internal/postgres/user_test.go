package postgres_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/core"
	"github.com/jacob-ian/jacobianmatthews.com/backend/internal/postgres"
)

type userSuite struct {
	ExpectedQuery string
	Tests         []userTest
	Columns       []string
	ToValueSlice  func(u core.User) []driver.Value
}

type userTest struct {
	DriverError   error
	ExpectedError error
	Rows          [][]driver.Value
	Fn            func(ur *postgres.UserRepository) ([]core.User, error)
}

func runUserMockDbTests(t *testing.T, suite userSuite) {
	for i := range suite.Tests {
		test := suite.Tests[i]

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("%v: Unexpected error '%s' creating db mock", i, err.Error())
		}

		rows := mock.NewRows(suite.Columns)
		for j := range test.Rows {
			rows.AddRow(test.Rows[j]...)
		}

		ex := mock.ExpectQuery(suite.ExpectedQuery)
		if test.DriverError != nil {
			ex.WillReturnError(test.DriverError)
		} else {
			ex.WillReturnRows(rows)
		}

		ur := postgres.NewUserRepository(context.Background(), db)
		values, tErr := test.Fn(ur)

		if !errors.Is(tErr, test.ExpectedError) {
			t.Errorf("%v: Unexpected error got '%v' want '%v'", i, tErr, test.ExpectedError)
		}

		if test.DriverError != nil && tErr == nil {
			t.Errorf("%v: Expected error to not be nil", i)
		}

		for j := range test.Rows {
			row := test.Rows[j]
			value := suite.ToValueSlice(values[j])
			for k := range row {
				if got, want := row[k], value[k]; got != want {
					t.Errorf("%v: Unexpected value got %v want %v", i, got, want)
				}
			}
		}

		expMet := mock.ExpectationsWereMet()
		if expMet != nil {
			t.Errorf("%v: Expectations were not met: %v", i, expMet)
		}
	}
}

func TestUserRepository_FindById(t *testing.T) {
	runUserMockDbTests(t, userSuite{
		ExpectedQuery: `SELECT \* FROM users`,
		Columns:       []string{"id", "name", "email", "email_verified", "image_url", "created_at", "updated_at", "deleted_at"},
		ToValueSlice: func(value core.User) []driver.Value {
			return []driver.Value{
				value.Id,
				value.Name,
				value.Email,
				value.EmailVerified,
				value.ImageUrl,
				value.CreatedAt,
				value.UpdatedAt,
				value.DeletedAt,
			}
		},
		Tests: []userTest{
			{
				DriverError:   nil,
				ExpectedError: nil,
				Rows: [][]driver.Value{
					{
						uuid.Must(uuid.NewRandom()).String(),
						"Testy test",
						"emaily@email.com",
						true,
						"https://google.com/logo.png",
						time.Now(),
						time.Now(),
						time.Time{},
					},
				},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					id := uuid.Must(uuid.NewRandom()).String()
					user, err := ur.FindById(context.Background(), id)
					return []core.User{user}, err
				},
			},
			{
				DriverError:   sql.ErrNoRows,
				ExpectedError: core.NewError(core.NotFoundError, "User not found"),
				Rows:          [][]driver.Value{},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					user, err := ur.FindById(context.Background(), uuid.Must(uuid.NewRandom()).String())
					return []core.User{user}, err
				},
			},
			{
				DriverError:   sql.ErrConnDone,
				ExpectedError: core.NewError(core.InternalError, "Could not get user"),
				Rows:          [][]driver.Value{},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					user, err := ur.FindById(context.Background(), uuid.Must(uuid.NewRandom()).String())
					return []core.User{user}, err
				},
			},
		},
	})
}

func TestUserRepository_FindAll(t *testing.T) {
	runUserMockDbTests(t, userSuite{
		ExpectedQuery: `SELECT \* FROM users`,
		Columns:       []string{"id", "name", "email", "email_verified", "image_url", "created_at", "updated_at", "deleted_at"},
		ToValueSlice: func(value core.User) []driver.Value {
			return []driver.Value{
				value.Id,
				value.Name,
				value.Email,
				value.EmailVerified,
				value.ImageUrl,
				value.CreatedAt,
				value.UpdatedAt,
				value.DeletedAt,
			}
		},
		Tests: []userTest{
			{
				DriverError:   nil,
				ExpectedError: nil,
				Rows: [][]driver.Value{
					{
						uuid.Must(uuid.NewRandom()).String(),
						"Testy test",
						"emaily@email.com",
						true,
						"https://google.com/logo.png",
						time.Now(),
						time.Now(),
						time.Time{},
					},
					{
						uuid.Must(uuid.NewRandom()).String(),
						"Usery user",
						"usery@email.com",
						false,
						"https://google.com/logo1.png",
						time.Now(),
						time.Now(),
						time.Time{},
					},
				},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					return ur.FindAll(context.Background())
				},
			},
			{
				DriverError:   nil,
				ExpectedError: nil,
				Rows:          [][]driver.Value{},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					return ur.FindAll(context.Background())
				},
			},
			{
				DriverError:   sql.ErrConnDone,
				ExpectedError: core.NewError(core.InternalError, "Could not get users"),
				Rows:          [][]driver.Value{},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					return ur.FindAll(context.Background())
				},
			},
		},
	})
}

func TestUserRepository_Create(t *testing.T) {
	newUser := core.NewUser{
		Id:            uuid.Must(uuid.NewRandom()).String(),
		Name:          "New User",
		Email:         "new@email.com",
		EmailVerified: false,
	}
	runUserMockDbTests(t, userSuite{
		ExpectedQuery: `INSERT INTO users`,
		Columns:       []string{"id", "name", "email", "email_verified", "image_url", "created_at", "updated_at", "deleted_at"},
		ToValueSlice: func(value core.User) []driver.Value {
			return []driver.Value{
				value.Id,
				value.Name,
				value.Email,
				value.EmailVerified,
				value.ImageUrl,
				value.CreatedAt,
				value.UpdatedAt,
				value.DeletedAt,
			}
		},
		Tests: []userTest{
			{
				DriverError:   nil,
				ExpectedError: nil,
				Rows: [][]driver.Value{
					{
						newUser.Id,
						newUser.Name,
						newUser.Email,
						newUser.EmailVerified,
						"",
						time.Now(),
						time.Now(),
						time.Time{},
					},
				},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					user, err := ur.Create(context.Background(), newUser)
					return []core.User{user}, err
				},
			},
			{
				DriverError:   sql.ErrNoRows,
				ExpectedError: core.NewError(core.InternalError, "Could not create user"),
				Rows:          [][]driver.Value{},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					user, err := ur.Create(context.Background(), newUser)
					return []core.User{user}, err
				},
			},
			{
				DriverError:   sql.ErrConnDone,
				ExpectedError: core.NewError(core.InternalError, "Could not create user"),
				Rows:          [][]driver.Value{},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					user, err := ur.Create(context.Background(), newUser)
					return []core.User{user}, err
				},
			},
		},
	})
}

func TestUserRepository_Update(t *testing.T) {
	userUpdate := core.User{
		Id:            uuid.Must(uuid.NewRandom()).String(),
		Name:          "Updated User",
		Email:         "updated@email.com",
		EmailVerified: true,
		ImageUrl:      "fake",
		CreatedAt:     time.Date(2022, time.September, 13, 12, 0, 0, 0, time.Local),
		UpdatedAt:     time.Now(),
	}
	runUserMockDbTests(t, userSuite{
		ExpectedQuery: `UPDATE users`,
		Columns:       []string{"id", "name", "email", "email_verified", "image_url", "created_at", "updated_at", "deleted_at"},
		ToValueSlice: func(value core.User) []driver.Value {
			return []driver.Value{
				value.Id,
				value.Name,
				value.Email,
				value.EmailVerified,
				value.ImageUrl,
				value.CreatedAt,
				value.UpdatedAt,
				value.DeletedAt,
			}
		},
		Tests: []userTest{
			{
				DriverError:   nil,
				ExpectedError: nil,
				Rows: [][]driver.Value{
					{
						userUpdate.Id,
						userUpdate.Name,
						userUpdate.Email,
						userUpdate.EmailVerified,
						userUpdate.ImageUrl,
						userUpdate.CreatedAt,
						userUpdate.UpdatedAt,
						userUpdate.DeletedAt,
					},
				},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					user, err := ur.Update(context.Background(), userUpdate)
					return []core.User{user}, err
				},
			},
			{
				DriverError:   sql.ErrNoRows,
				ExpectedError: core.NewError(core.NotFoundError, "User not found"),
				Rows:          [][]driver.Value{},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					user, err := ur.Update(context.Background(), userUpdate)
					return []core.User{user}, err
				},
			},
			{
				DriverError:   sql.ErrConnDone,
				ExpectedError: core.NewError(core.InternalError, "Could not update user"),
				Rows:          [][]driver.Value{},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					user, err := ur.Update(context.Background(), userUpdate)
					return []core.User{user}, err
				},
			},
		},
	})
}

func TestUserRepository_Delete(t *testing.T) {
	runUserMockDbTests(t, userSuite{
		ExpectedQuery: `UPDATE users`,
		Columns:       []string{"id", "name", "email", "email_verified", "image_url", "created_at", "updated_at", "deleted_at"},
		ToValueSlice: func(value core.User) []driver.Value {
			return []driver.Value{
				value.Id,
				value.Name,
				value.Email,
				value.EmailVerified,
				value.ImageUrl,
				value.CreatedAt,
				value.UpdatedAt,
				value.DeletedAt,
			}
		},
		Tests: []userTest{
			{
				DriverError:   nil,
				ExpectedError: nil,
				Rows: [][]driver.Value{
					{
						"",
						"",
						"",
						false,
						"",
						time.Time{},
						time.Time{},
						time.Time{},
					},
				},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					id := uuid.Must(uuid.NewRandom()).String()
					err := ur.Delete(context.Background(), id)
					return []core.User{core.User{}}, err
				},
			},
			{
				DriverError:   sql.ErrNoRows,
				ExpectedError: core.NewError(core.NotFoundError, "User not found"),
				Rows:          [][]driver.Value{},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					id := uuid.Must(uuid.NewRandom()).String()
					err := ur.Delete(context.Background(), id)
					return []core.User{core.User{}}, err
				},
			},
			{
				DriverError:   sql.ErrConnDone,
				ExpectedError: core.NewError(core.InternalError, "Could not delete user"),
				Rows:          [][]driver.Value{},
				Fn: func(ur *postgres.UserRepository) ([]core.User, error) {
					id := uuid.Must(uuid.NewRandom()).String()
					err := ur.Delete(context.Background(), id)
					return []core.User{core.User{}}, err
				},
			},
		},
	})
}
