package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RoleMenuModel = (*customRoleMenuModel)(nil)

type (
	// RoleMenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleMenuModel.
	RoleMenuModel interface {
		roleMenuModel
		FindMenuIdsByRoleId(ctx context.Context, roleId uint64) ([]uint64, error)
		withSession(session sqlx.Session) RoleMenuModel
	}

	customRoleMenuModel struct {
		*defaultRoleMenuModel
	}
)

// NewRoleMenuModel returns a model for the database table.
func NewRoleMenuModel(conn sqlx.SqlConn) RoleMenuModel {
	return &customRoleMenuModel{
		defaultRoleMenuModel: newRoleMenuModel(conn),
	}
}

func (m *customRoleMenuModel) withSession(session sqlx.Session) RoleMenuModel {
	return NewRoleMenuModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customRoleMenuModel) FindMenuIdsByRoleId(ctx context.Context, roleId uint64) ([]uint64, error) {
	query := fmt.Sprintf("select `menu_id` from %s where `role_id` = ?", m.table)
	var ids []uint64
	if err := m.conn.QueryRowsCtx(ctx, &ids, query, roleId); err != nil {
		return nil, err
	}
	return ids, nil
}
