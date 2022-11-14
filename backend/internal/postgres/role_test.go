package postgres_test

import (
	"context"
	"database/sql"
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
		Rows          []core.Role
		Error         error
	}

	tests := []test{
		{
			RoleName: "Namey Namer",
			Columns:  []string{"id", "name", "created_at", "deleted_at"},
			Rows: []core.Role{{
				Id:        uuid.Must(uuid.NewRandom()),
				Name:      "Namey Namer",
				CreatedAt: time.Now(),
				DeletedAt: time.Time{},
			}},
			ExpectedQuery: "SELECT * FROM roles",
			Error:         nil,
		},
		{
			RoleName:      "Namey Namer",
			Columns:       []string{"id", "name", "created_at", "deleted_at"},
			Rows:          []core.Role{},
			ExpectedQuery: "SELECT * FROM roles",
			Error:         sql.ErrNoRows,
		},
	}

	for i := range tests {
		test := tests[i]
		rows := mock.NewRows(test.Columns)
		for j := range test.Rows {
			rows.AddRow(test.Rows[j])
		}
		mock.ExpectBegin()
		ex := mock.ExpectQuery(test.ExpectedQuery).WithArgs(test.RoleName)
		if test.Error != nil {
			ex.WillReturnError(test.Error)
		} else {
			ex.WillReturnRows(rows)
		}

		r := postgres.NewRoleRepository(context.Background(), db)
		role, tErr := r.FindByName(context.Background(), test.RoleName)
		if tErr != nil && test.Error == nil {
			t.Errorf("Unexpected error '%s' when running test", tErr.Error())
		}
		if tErr == nil && test.Error != nil {
			t.Errorf("Expected error '%s' when running test", test.Error)
		}
		if role.Name != test.RoleName {
			t.Errorf("Unexpected name got '%s' want '%s'", role.Name, test.RoleName)
		}
	}

}
