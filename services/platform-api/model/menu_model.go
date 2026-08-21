package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MenuModel = (*customMenuModel)(nil)

type (
	// MenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMenuModel.
	MenuModel interface {
		menuModel
		FindByIds(ctx context.Context, ids []uint64) ([]*Menu, error)
		FindByRoleId(ctx context.Context, roleId uint64) ([]*Menu, error)
		withSession(session sqlx.Session) MenuModel
	}

	customMenuModel struct {
		*defaultMenuModel
	}
)

// NewMenuModel returns a model for the database table.
func NewMenuModel(conn sqlx.SqlConn) MenuModel {
	return &customMenuModel{
		defaultMenuModel: newMenuModel(conn),
	}
}

func (m *customMenuModel) withSession(session sqlx.Session) MenuModel {
	return NewMenuModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customMenuModel) FindByIds(ctx context.Context, ids []uint64) ([]*Menu, error) {
	if len(ids) == 0 {
		return []*Menu{}, nil
	}

	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]
	query := fmt.Sprintf("select %s from %s where `id` in (%s) order by `id`", menuRows, m.table, placeholders)

	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	var list []*Menu
	if err := m.conn.QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customMenuModel) FindByRoleId(ctx context.Context, roleId uint64) ([]*Menu, error) {
	query := fmt.Sprintf(
		"select %s from %s where `id` in (select `menu_id` from `role_menu` where `role_id` = ?) order by `id`",
		menuRows, m.table,
	)
	var list []*Menu
	if err := m.conn.QueryRowsCtx(ctx, &list, query, roleId); err != nil {
		return nil, err
	}
	return list, nil
}
