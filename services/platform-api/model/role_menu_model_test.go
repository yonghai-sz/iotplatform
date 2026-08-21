package model

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func newMockRoleMenuModel(t *testing.T) (RoleMenuModel, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return NewRoleMenuModel(sqlx.NewSqlConnFromDB(db)), mock
}

func TestRoleMenuModel_FindMenuIdsByRoleId(t *testing.T) {
	ast := assert.New(t)
	m, mock := newMockRoleMenuModel(t)
	ctx := context.Background()

	mock.ExpectQuery("select `menu_id` from `role_menu` where `role_id` = \\?").
		WithArgs(uint64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"menu_id"}).AddRow(uint64(2)).AddRow(uint64(54)))

	ids, err := m.FindMenuIdsByRoleId(ctx, 2)
	ast.NoError(err)
	ast.Equal([]uint64{2, 54}, ids)
}

func TestRoleMenuModel_FindOneByRoleIdMenuId(t *testing.T) {
	ast := assert.New(t)
	m, mock := newMockRoleMenuModel(t)
	ctx := context.Background()

	mock.ExpectQuery("select .+ from `role_menu` where `role_id` = \\? and `menu_id` = \\?").
		WithArgs(uint64(2), uint64(54)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role_id", "menu_id"}).
			AddRow(uint64(1), uint64(2), uint64(54)))

	got, err := m.FindOneByRoleIdMenuId(ctx, 2, 54)
	ast.NoError(err)
	ast.Equal(uint64(2), got.RoleId)
	ast.Equal(uint64(54), got.MenuId)
}
