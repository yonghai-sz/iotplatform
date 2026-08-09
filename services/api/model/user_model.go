package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		FindPage(ctx context.Context, tenantId uint64, username string, pageIndex, pageSize int64) (int64, []*User, error)
		withSession(session sqlx.Session) UserModel
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

func (m *customUserModel) withSession(session sqlx.Session) UserModel {
	return NewUserModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	now := time.Now()
	data.CreatedAt = sql.NullTime{Time: now, Valid: true}
	data.UpdatedAt = sql.NullTime{Time: now, Valid: true}
	data.DeletedAt = sql.NullTime{}
	if data.Enable == "" {
		data.Enable = "Enable"
	}
	query := fmt.Sprintf("insert into %s (`created_at`, `updated_at`, `username`, `password`, `salt`, `enable`, `role_id`, `tenant_id`) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table)
	return m.conn.ExecCtx(ctx, query, data.CreatedAt, data.UpdatedAt, data.Username, data.Password, data.Salt, data.Enable, data.RoleId, data.TenantId)
}

func (m *customUserModel) FindOne(ctx context.Context, id uint64) (*User, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `deleted_at` is null limit 1", userRows, m.table)
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserModel) FindOneByUsername(ctx context.Context, username string) (*User, error) {
	query := fmt.Sprintf("select %s from %s where `username` = ? and `deleted_at` is null limit 1", userRows, m.table)
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, username)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserModel) Update(ctx context.Context, data *User) error {
	data.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	query := fmt.Sprintf("update %s set `updated_at` = ?, `username` = ?, `password` = ?, `salt` = ?, `enable` = ?, `role_id` = ?, `tenant_id` = ? where `id` = ? and `deleted_at` is null", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.UpdatedAt, data.Username, data.Password, data.Salt, data.Enable, data.RoleId, data.TenantId, data.Id)
	return err
}

func (m *customUserModel) Delete(ctx context.Context, id uint64) error {
	now := time.Now()
	query := fmt.Sprintf("update %s set `deleted_at` = ?, `updated_at` = ? where `id` = ? and `deleted_at` is null", m.table)
	result, err := m.conn.ExecCtx(ctx, query, now, now, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func (m *customUserModel) FindPage(ctx context.Context, tenantId uint64, username string, pageIndex, pageSize int64) (int64, []*User, error) {
	var (
		conds []string
		args  []any
	)
	conds = append(conds, "`deleted_at` is null")
	if tenantId > 0 {
		conds = append(conds, "`tenant_id` = ?")
		args = append(args, tenantId)
	}
	if username != "" {
		conds = append(conds, "`username` like ?")
		args = append(args, "%"+username+"%")
	}
	where := " where " + strings.Join(conds, " and ")

	var total int64
	countQuery := fmt.Sprintf("select count(*) from %s%s", m.table, where)
	if err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, err
	}

	offset := (pageIndex - 1) * pageSize
	listQuery := fmt.Sprintf("select %s from %s%s order by `updated_at` desc limit ? offset ?", userRows, m.table, where)
	listArgs := append(append([]any{}, args...), pageSize, offset)

	var list []*User
	if err := m.conn.QueryRowsCtx(ctx, &list, listQuery, listArgs...); err != nil {
		return 0, nil, err
	}
	return total, list, nil
}
