package model

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func newMockUserModel(t *testing.T) (UserModel, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return NewUserModel(sqlx.NewSqlConnFromDB(db)), mock
}

func TestUserModel_InsertAndFindOneByUsername(t *testing.T) {
	ast := assert.New(t)
	m, mock := newMockUserModel(t)
	ctx := context.Background()

	mock.ExpectExec("insert into `user`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "alice", "hash", "salt", "Enable", uint64(2), uint64(1)).
		WillReturnResult(sqlmock.NewResult(8, 1))

	result, err := m.Insert(ctx, &User{
		Username: "alice",
		Password: "hash",
		Salt:     "salt",
		Enable:   "Enable",
		RoleId:   2,
		TenantId: 1,
	})
	ast.NoError(err)
	id, err := result.LastInsertId()
	ast.NoError(err)
	ast.Equal(int64(8), id)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "password", "salt", "enable", "role_id", "tenant_id"}).
		AddRow(uint64(8), now, now, nil, "alice", "hash", "salt", "Enable", uint64(2), uint64(1))
	mock.ExpectQuery("select .+ from `user` where `username` = \\? and `deleted_at` is null").
		WithArgs("alice").
		WillReturnRows(rows)

	got, err := m.FindOneByUsername(ctx, "alice")
	ast.NoError(err)
	ast.Equal("alice", got.Username)
	ast.Equal(uint64(2), got.RoleId)
}

func TestUserModel_FindPage(t *testing.T) {
	ast := assert.New(t)
	m, mock := newMockUserModel(t)
	ctx := context.Background()

	mock.ExpectQuery("select count\\(\\*\\) from `user` where `deleted_at` is null and `tenant_id` = \\? and `username` like \\?").
		WithArgs(uint64(1), "%alice%").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))

	now := time.Now()
	mock.ExpectQuery("select .+ from `user` where `deleted_at` is null and `tenant_id` = \\? and `username` like \\? order by `updated_at` desc limit \\? offset \\?").
		WithArgs(uint64(1), "%alice%", int64(20), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "password", "salt", "enable", "role_id", "tenant_id"}).
			AddRow(uint64(1), now, now, nil, "alice", "hash", "salt", "Enable", uint64(2), uint64(1)))

	total, list, err := m.FindPage(ctx, 1, "alice", 1, 20)
	ast.NoError(err)
	ast.Equal(int64(1), total)
	ast.Len(list, 1)
	ast.Equal("alice", list[0].Username)
}

func TestUserModel_SoftDelete(t *testing.T) {
	ast := assert.New(t)
	m, mock := newMockUserModel(t)
	ctx := context.Background()

	mock.ExpectExec(fmt.Sprintf("update %s set `deleted_at`", "`user`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), uint64(3)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	ast.NoError(m.Delete(ctx, 3))

	mock.ExpectExec(fmt.Sprintf("update %s set `deleted_at`", "`user`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), uint64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	ast.ErrorIs(m.Delete(ctx, 99), ErrNotFound)
}
