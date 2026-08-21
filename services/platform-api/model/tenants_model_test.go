package model

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func newMockTenantsModel(t *testing.T) (TenantsModel, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return NewTenantsModel(sqlx.NewSqlConnFromDB(db)), mock
}

func TestTenantsModel_InsertAndFindOne(t *testing.T) {
	ast := assert.New(t)
	m, mock := newMockTenantsModel(t)
	ctx := context.Background()

	mock.ExpectExec("insert into `tenants`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "acme", "a@acme.com").
		WillReturnResult(sqlmock.NewResult(7, 1))

	result, err := m.Insert(ctx, &Tenants{
		TenantName: "acme",
		Email:      sql.NullString{String: "a@acme.com", Valid: true},
	})
	ast.NoError(err)
	id, err := result.LastInsertId()
	ast.NoError(err)
	ast.Equal(int64(7), id)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "tenant_name", "email"}).
		AddRow(uint64(7), now, now, nil, "acme", "a@acme.com")
	mock.ExpectQuery("select .+ from `tenants` where `id` = \\? and `deleted_at` is null").
		WithArgs(uint64(7)).
		WillReturnRows(rows)

	got, err := m.FindOne(ctx, 7)
	ast.NoError(err)
	ast.Equal("acme", got.TenantName)
	ast.Equal("a@acme.com", got.Email.String)
}

func TestTenantsModel_FindPage(t *testing.T) {
	ast := assert.New(t)
	m, mock := newMockTenantsModel(t)
	ctx := context.Background()

	mock.ExpectQuery("select count\\(\\*\\) from `tenants` where `deleted_at` is null and `tenant_name` like \\?").
		WithArgs("%acme%").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))

	now := time.Now()
	mock.ExpectQuery("select .+ from `tenants` where `deleted_at` is null and `tenant_name` like \\? order by `updated_at` desc limit \\? offset \\?").
		WithArgs("%acme%", int64(20), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "tenant_name", "email"}).
			AddRow(uint64(1), now, now, nil, "acme", "a@acme.com"))

	total, list, err := m.FindPage(ctx, "acme", 1, 20)
	ast.NoError(err)
	ast.Equal(int64(1), total)
	ast.Len(list, 1)
	ast.Equal("acme", list[0].TenantName)
}

func TestTenantsModel_SoftDelete(t *testing.T) {
	ast := assert.New(t)
	m, mock := newMockTenantsModel(t)
	ctx := context.Background()

	mock.ExpectExec(fmt.Sprintf("update %s set `deleted_at`", "`tenants`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), uint64(3)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ast.NoError(m.Delete(ctx, 3))

	mock.ExpectExec(fmt.Sprintf("update %s set `deleted_at`", "`tenants`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), uint64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	ast.ErrorIs(m.Delete(ctx, 99), ErrNotFound)
}
