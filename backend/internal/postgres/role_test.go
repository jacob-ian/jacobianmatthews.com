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

func TestRoleRepository_FindByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Unexpected error '%s' when creating stub database", err.Error())
	}
	defer db.Close()

	type test struct {
		ExpectedQuery string
		RoleName      string
		Columns       []string
		Rows          [][]driver.Value
		DriverError   error
		ExpectedError error
	}

	tests := []test{
		{
			RoleName: "Namey Namer",
			Columns:  []string{"id", "name", "created_at", "deleted_at"},
			Rows: [][]driver.Value{
				{uuid.Must(uuid.NewRandom()), "Namey Namer", time.Now(), time.Time{}},
			},
			ExpectedQuery: `SELECT \* FROM roles`,
			DriverError:   nil,
			ExpectedError: nil,
		},
		{
			RoleName:      "Namey Namer",
			Columns:       []string{"id", "name", "created_at", "deleted_at"},
			Rows:          [][]driver.Value{},
			ExpectedQuery: `SELECT \* FROM roles`,
			DriverError:   sql.ErrNoRows,
			ExpectedError: core.NewError(core.NotFoundError, "Role not found"),
		},
	}

	for i := range tests {
		test := tests[i]

		rows := mock.NewRows(test.Columns)
		for j := range test.Rows {
			values := test.Rows[j]
			rows.AddRow(values...)
		}

		ex := mock.ExpectQuery(test.ExpectedQuery).WithArgs(test.RoleName)
		if test.DriverError != nil {
			ex.WillReturnError(test.DriverError)
		} else {
			ex.WillReturnRows(rows)
		}

		r := postgres.NewRoleRepository(context.Background(), db)
		role, tErr := r.FindByName(context.Background(), test.RoleName)

		if test.DriverError != nil {
			if !errors.Is(tErr, test.ExpectedError) {
				t.Errorf("Expected error got '%s' want '%s'", tErr.Error(), test.DriverError.Error())
			}
			return
		}

		if tErr != nil {
			t.Errorf("Unexpected error '%s'", tErr.Error())
		}

		if role.Name != test.RoleName {
			t.Errorf("Unexpected name got '%s' want '%s'", role.Name, test.RoleName)
		}
	}
}

func TestRoleRepository_FindAll(t *testing.T) {
	type test struct {
		ExpectedQuery string
		Columns       []string
		Rows          [][]driver.Value
		DriverError   error
		ExpectedError error
	}

	tests := []test{
		{
			ExpectedQuery: `SELECT \* FROM roles`,
			Columns:       []string{"id", "name", "created_at", "deleted_at"},
			Rows: [][]driver.Value{
				{
					uuid.Must(uuid.NewRandom()),
					"Admin",
					time.Now(),
					time.Time{},
				},
				{
					uuid.Must(uuid.NewRandom()),
					"Author",
					time.Now(),
					time.Time{},
				},
				{
					uuid.Must(uuid.NewRandom()),
					"User",
					time.Now(),
					time.Time{},
				},
			},
			DriverError:   nil,
			ExpectedError: nil,
		},
		{
			ExpectedQuery: `SELECT \* FROM roles`,
			Columns:       []string{"id", "name", "created_at", "deleted_at"},
			Rows:          [][]driver.Value{},
			DriverError:   nil,
			ExpectedError: nil,
		},
		{
			ExpectedQuery: `SELECT \* FROM roles`,
			Columns:       []string{"id", "name", "created_at", "deleted_at"},
			Rows:          [][]driver.Value{},
			DriverError:   sql.ErrConnDone,
			ExpectedError: core.NewError(core.InternalError, "Could not get roles"),
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Unexpected error '%s' creating db mock", err.Error())
	}

	for i := range tests {
		test := tests[i]

		rows := mock.NewRows(test.Columns)
		for j := range test.Rows {
			rows.AddRow(test.Rows[j]...)
		}

		ex := mock.ExpectQuery(test.ExpectedQuery)
		if test.DriverError != nil {
			ex.WillReturnError(test.DriverError)
		} else {
			ex.WillReturnRows(rows)
		}

		rr := postgres.NewRoleRepository(context.Background(), db)
		roles, tErr := rr.FindAll(context.Background())
		if !errors.Is(tErr, test.ExpectedError) {
			t.Errorf("Unexpected error got '%v' want '%v'", tErr, test.ExpectedError)
		}

		if test.DriverError != nil && tErr == nil {
			t.Errorf("Expected error to not be nil")
		}

		if len(roles) != len(test.Rows) {
			t.Errorf("Unexpected roles length got '%v' want '%v'", len(roles), len(test.Rows))
		}

		for k := range test.Rows {
			row := test.Rows[k]
			role := roles[k]

			if role.Id != row[0] || role.Name != row[1] || role.CreatedAt != row[2] || role.DeletedAt != row[3] {
				t.Errorf("Unexpected row values got '%v' want '%v'", row, role)
			}
		}

		expMet := mock.ExpectationsWereMet()
		if expMet != nil {
			t.Errorf("Expectations were not met: %v", expMet)
		}
	}
}

func TestRoleRepository_Create(t *testing.T) {
	type test struct {
		RoleName      string
		ExpectedQuery string
		Columns       []string
		Result        []driver.Value
		DriverError   error
		ExpectedError error
	}

	tests := []test{
		{
			RoleName:      "Admin",
			ExpectedQuery: `INSERT INTO roles`,
			Columns:       []string{"id", "name", "created_at", "deleted_at"},
			Result: []driver.Value{
				uuid.Must(uuid.NewRandom()),
				"Admin",
				time.Now(),
				time.Time{},
			},
			DriverError:   nil,
			ExpectedError: nil,
		},
		{
			RoleName:      "Admin",
			ExpectedQuery: `INSERT INTO roles`,
			Columns:       []string{"id", "name", "created_at", "deleted_at"},
			Result:        nil,
			DriverError:   sql.ErrNoRows,
			ExpectedError: core.NewError(core.InternalError, "Could not create role"),
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Unexpected error '%s' creating db mock", err.Error())
	}

	for i := range tests {
		test := tests[i]

		rows := mock.NewRows(test.Columns)
		if test.Result != nil {
			rows.AddRow(test.Result...)
		}

		ex := mock.ExpectQuery(test.ExpectedQuery).WithArgs(test.RoleName)
		if test.DriverError != nil {
			ex.WillReturnError(test.DriverError)
		} else {
			ex.WillReturnRows(rows)
		}

		rr := postgres.NewRoleRepository(context.Background(), db)
		role, tErr := rr.Create(context.Background(), test.RoleName)
		if !errors.Is(tErr, test.ExpectedError) {
			t.Errorf("Unexpected error got '%v' want '%v'", tErr, test.ExpectedError)
		}

		if test.DriverError != nil && tErr == nil {
			t.Errorf("Expected error to not be nil")
		}

		if test.Result != nil && role.Name != test.RoleName {
			t.Errorf("Test %v: Unexpected role name got '%v' want '%v'", i, role.Name, test.RoleName)
		}

		expMet := mock.ExpectationsWereMet()
		if expMet != nil {
			t.Errorf("Expectations were not met: %v", expMet)
		}
	}
}

func TestRoleRepository_Delete(t *testing.T) {
	type test struct {
		RoleName      string
		ExpectedQuery string
		DriverError   error
		ExpectedError error
	}

	tests := []test{
		{
			RoleName:      "Admin",
			ExpectedQuery: `UPDATE roles`,
			DriverError:   nil,
			ExpectedError: nil,
		},
		{
			RoleName:      "Admin",
			ExpectedQuery: `UPDATE roles`,
			DriverError:   sql.ErrNoRows,
			ExpectedError: core.NewError(core.NotFoundError, "Role not found"),
		},
		{
			RoleName:      "Admin",
			ExpectedQuery: `UPDATE roles`,
			DriverError:   sql.ErrConnDone,
			ExpectedError: core.NewError(core.InternalError, "Could not delete role"),
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Unexpected error '%s' creating db mock", err.Error())
	}

	for i := range tests {
		test := tests[i]

		ex := mock.ExpectQuery(test.ExpectedQuery).WithArgs(test.RoleName)
		if test.DriverError != nil {
			ex.WillReturnError(test.DriverError)
		} else {
			ex.WillReturnRows(mock.NewRows([]string{}))
		}

		rr := postgres.NewRoleRepository(context.Background(), db)
		tErr := rr.Delete(context.Background(), test.RoleName)
		if !errors.Is(tErr, test.ExpectedError) {
			t.Errorf("Unexpected error got '%v' want '%v'", tErr, test.ExpectedError)
		}

		if test.DriverError != nil && tErr == nil {
			t.Errorf("Expected error to not be nil")
		}

		expMet := mock.ExpectationsWereMet()
		if expMet != nil {
			t.Errorf("Expectations were not met: %v", expMet)
		}
	}

}
