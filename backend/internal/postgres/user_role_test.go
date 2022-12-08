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

type userRoleSuite struct {
	ExpectedQuery string
	Tests         []userRoleTest
	Columns       []string
	ToValueSlice  func(u any) []driver.Value
}

type userRoleTest struct {
	DriverError   error
	ExpectedError error
	Rows          [][]driver.Value
	Fn            func(ur *postgres.UserRoleRepository) ([]any, error)
}

func runUserRoleMockDbTests(t *testing.T, suite userRoleSuite) {
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

		ur := postgres.NewUserRoleRepository(context.Background(), db)
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

func TestUserRoleRepository_FindRoleByUserId(t *testing.T) {
	runUserRoleMockDbTests(t, userRoleSuite{
		ExpectedQuery: `SELECT roles.\* FROM roles`,
		Columns:       []string{"id", "name", "created_at", "deleted_at"},
		ToValueSlice: func(value any) []driver.Value {
			role := value.(core.Role)
			return []driver.Value{
				role.Id,
				role.Name,
				role.CreatedAt,
				role.DeletedAt,
			}
		},
		Tests: []userRoleTest{
			{
				DriverError:   nil,
				ExpectedError: nil,
				Rows: [][]driver.Value{
					{
						uuid.Must(uuid.NewRandom()),
						"Admin",
						time.Now(),
						time.Time{},
					},
				},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					id := uuid.Must(uuid.NewRandom()).String()
					role, err := urr.FindRoleByUserId(context.Background(), id)
					return []any{role}, err
				},
			},
			{
				DriverError:   sql.ErrNoRows,
				ExpectedError: core.NewError(core.NotFoundError, "Role not found"),
				Rows:          [][]driver.Value{},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					id := uuid.Must(uuid.NewRandom()).String()
					role, err := urr.FindRoleByUserId(context.Background(), id)
					return []any{role}, err
				},
			},
			{
				DriverError:   sql.ErrConnDone,
				ExpectedError: core.NewError(core.InternalError, "Could not find role"),
				Rows:          [][]driver.Value{},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					id := uuid.Must(uuid.NewRandom()).String()
					role, err := urr.FindRoleByUserId(context.Background(), id)
					return []any{role}, err
				},
			},
		},
	})
}

func TestUserRoleRepository_FindById(t *testing.T) {
	id := uuid.Must(uuid.NewRandom())
	runUserRoleMockDbTests(t, userRoleSuite{
		ExpectedQuery: `SELECT \* FROM user_role`,
		Columns:       []string{"id", "user_id", "role_id", "created_at", "updated_at", "deleted_at"},
		ToValueSlice: func(value any) []driver.Value {
			userRole := value.(core.UserRole)
			return []driver.Value{
				userRole.Id,
				userRole.UserId,
				userRole.RoleId,
				userRole.CreatedAt,
				userRole.UpdatedAt,
				userRole.DeletedAt,
			}
		},
		Tests: []userRoleTest{
			{
				DriverError:   nil,
				ExpectedError: nil,
				Rows: [][]driver.Value{
					{
						id,
						uuid.Must(uuid.NewRandom()).String(),
						uuid.Must(uuid.NewRandom()).String(),
						time.Now(),
						time.Now(),
						time.Time{},
					},
				},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					ur, err := urr.FindById(context.Background(), id)
					return []any{ur}, err
				},
			},
			{
				DriverError:   sql.ErrNoRows,
				ExpectedError: core.NewError(core.NotFoundError, "User role not found"),
				Rows:          [][]driver.Value{},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					ur, err := urr.FindById(context.Background(), id)
					return []any{ur}, err
				},
			},
			{
				DriverError:   sql.ErrConnDone,
				ExpectedError: core.NewError(core.InternalError, "Could not find user role"),
				Rows:          [][]driver.Value{},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					ur, err := urr.FindById(context.Background(), id)
					return []any{ur}, err
				},
			},
		},
	})
}

func TestUserRoleRepository_Create(t *testing.T) {
	userId := uuid.Must(uuid.NewRandom()).String()
	roleId := uuid.Must(uuid.NewRandom())
	runUserRoleMockDbTests(t, userRoleSuite{
		ExpectedQuery: `INSERT INTO user_role`,
		Columns:       []string{"id", "user_id", "role_id", "created_at", "updated_at", "deleted_at"},
		ToValueSlice: func(value any) []driver.Value {
			userRole := value.(core.UserRole)
			return []driver.Value{
				userRole.Id,
				userRole.UserId,
				userRole.RoleId,
				userRole.CreatedAt,
				userRole.UpdatedAt,
				userRole.DeletedAt,
			}
		},
		Tests: []userRoleTest{
			{
				DriverError:   nil,
				ExpectedError: nil,
				Rows: [][]driver.Value{
					{
						uuid.Must(uuid.NewRandom()),
						userId,
						roleId.String(),
						time.Now(),
						time.Now(),
						time.Time{},
					},
				},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					ur, err := urr.Create(context.Background(), userId, roleId)
					return []any{ur}, err
				},
			},
			{
				DriverError:   sql.ErrNoRows,
				ExpectedError: core.NewError(core.InternalError, "Could not create user role"),
				Rows:          [][]driver.Value{},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					ur, err := urr.Create(context.Background(), userId, roleId)
					return []any{ur}, err
				},
			},
			{
				DriverError:   sql.ErrConnDone,
				ExpectedError: core.NewError(core.InternalError, "Could not create user role"),
				Rows:          [][]driver.Value{},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					ur, err := urr.Create(context.Background(), userId, roleId)
					return []any{ur}, err
				},
			},
		},
	})
}

func TestUserRoleRepository_UpdateByUserId(t *testing.T) {
	userId := uuid.Must(uuid.NewRandom()).String()
	roleId := uuid.Must(uuid.NewRandom())
	runUserRoleMockDbTests(t, userRoleSuite{
		ExpectedQuery: `UPDATE user_role`,
		Columns:       []string{"id", "user_id", "role_id", "created_at", "updated_at", "deleted_at"},
		ToValueSlice: func(value any) []driver.Value {
			userRole := value.(core.UserRole)
			return []driver.Value{
				userRole.Id,
				userRole.UserId,
				userRole.RoleId,
				userRole.CreatedAt,
				userRole.UpdatedAt,
				userRole.DeletedAt,
			}
		},
		Tests: []userRoleTest{
			{
				DriverError:   nil,
				ExpectedError: nil,
				Rows: [][]driver.Value{
					{
						uuid.Must(uuid.NewRandom()),
						userId,
						roleId.String(),
						time.Now(),
						time.Now(),
						time.Time{},
					},
				},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					ur, err := urr.UpdateByUserId(context.Background(), userId, roleId)
					return []any{ur}, err
				},
			},
			{
				DriverError:   sql.ErrNoRows,
				ExpectedError: core.NewError(core.NotFoundError, "User role not found"),
				Rows:          [][]driver.Value{},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					ur, err := urr.UpdateByUserId(context.Background(), userId, roleId)
					return []any{ur}, err
				},
			},
			{
				DriverError:   sql.ErrConnDone,
				ExpectedError: core.NewError(core.InternalError, "Could not update user role"),
				Rows:          [][]driver.Value{},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					ur, err := urr.UpdateByUserId(context.Background(), userId, roleId)
					return []any{ur}, err
				},
			},
		},
	})
}

func TestUserRoleRepository_DeleteByUserId(t *testing.T) {
	userId := uuid.Must(uuid.NewRandom()).String()
	runUserRoleMockDbTests(t, userRoleSuite{
		ExpectedQuery: `UPDATE user_role`,
		Columns:       []string{"id", "user_id", "role_id", "created_at", "updated_at", "deleted_at"},
		ToValueSlice: func(value any) []driver.Value {
			userRole := value.(core.UserRole)
			return []driver.Value{
				userRole.Id,
				userRole.UserId,
				userRole.RoleId,
				userRole.CreatedAt,
				userRole.UpdatedAt,
				userRole.DeletedAt,
			}
		},
		Tests: []userRoleTest{
			{
				DriverError:   nil,
				ExpectedError: nil,
				Rows:          [][]driver.Value{},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					err := urr.DeleteByUserId(context.Background(), userId)
					return []any{core.UserRole{}}, err
				},
			},
			{
				DriverError:   sql.ErrNoRows,
				ExpectedError: core.NewError(core.NotFoundError, "User role not found"),
				Rows:          [][]driver.Value{},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					err := urr.DeleteByUserId(context.Background(), userId)
					return []any{core.UserRole{}}, err
				},
			},
			{
				DriverError:   sql.ErrConnDone,
				ExpectedError: core.NewError(core.InternalError, "Could not delete user role"),
				Rows:          [][]driver.Value{},
				Fn: func(urr *postgres.UserRoleRepository) ([]any, error) {
					err := urr.DeleteByUserId(context.Background(), userId)
					return []any{core.UserRole{}}, err
				},
			},
		},
	})
}
