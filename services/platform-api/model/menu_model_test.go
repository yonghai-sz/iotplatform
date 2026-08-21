package model

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func newMockMenuModel(t *testing.T) (MenuModel, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return NewMenuModel(sqlx.NewSqlConnFromDB(db)), mock
}

func TestMenuModel_FindByIds(t *testing.T) {
	ast := assert.New(t)
	m, mock := newMockMenuModel(t)
	ctx := context.Background()

	list, err := m.FindByIds(ctx, nil)
	ast.NoError(err)
	ast.Empty(list)

	mock.ExpectQuery("select .+ from `menu` where `id` in \\(\\?,\\?\\) order by `id`").
		WithArgs(uint64(1), uint64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "parent_id", "menu_key", "title", "has_child"}).
			AddRow(uint64(1), uint64(0), "system", "系统管理", "Y").
			AddRow(uint64(2), uint64(1), "user", "用户管理", "N"))

	list, err = m.FindByIds(ctx, []uint64{1, 2})
	ast.NoError(err)
	ast.Len(list, 2)
	ast.Equal("system", list[0].MenuKey)
	ast.Equal("user", list[1].MenuKey)
}

func TestMenuModel_FindByRoleId(t *testing.T) {
	ast := assert.New(t)
	m, mock := newMockMenuModel(t)
	ctx := context.Background()

	mock.ExpectQuery("select .+ from `menu` where `id` in \\(select `menu_id` from `role_menu` where `role_id` = \\?\\) order by `id`").
		WithArgs(uint64(3)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "parent_id", "menu_key", "title", "has_child"}).
			AddRow(uint64(54), uint64(0), "videoCenter", "视频中心", "N"))

	list, err := m.FindByRoleId(ctx, 3)
	ast.NoError(err)
	ast.Len(list, 1)
	ast.Equal(uint64(54), list[0].Id)
	ast.Equal("videoCenter", list[0].MenuKey)
}
