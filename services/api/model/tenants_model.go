package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TenantsModel = (*customTenantsModel)(nil)

type (
	// TenantsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTenantsModel.
	TenantsModel interface {
		tenantsModel
		FindPage(ctx context.Context, tenantName string, pageIndex, pageSize int64) (int64, []*Tenants, error)
		withSession(session sqlx.Session) TenantsModel
	}

	customTenantsModel struct {
		*defaultTenantsModel
	}
)

// NewTenantsModel returns a model for the database table.
func NewTenantsModel(conn sqlx.SqlConn) TenantsModel {
	return &customTenantsModel{
		defaultTenantsModel: newTenantsModel(conn),
	}
}

func (m *customTenantsModel) withSession(session sqlx.Session) TenantsModel {
	return NewTenantsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customTenantsModel) Insert(ctx context.Context, data *Tenants) (sql.Result, error) {
	now := time.Now()
	data.CreatedAt = sql.NullTime{Time: now, Valid: true}
	data.UpdatedAt = sql.NullTime{Time: now, Valid: true}
	data.DeletedAt = sql.NullTime{}
	query := fmt.Sprintf("insert into %s (`created_at`, `updated_at`, `tenant_name`, `email`) values (?, ?, ?, ?)", m.table)
	return m.conn.ExecCtx(ctx, query, data.CreatedAt, data.UpdatedAt, data.TenantName, data.Email)
}

func (m *customTenantsModel) FindOne(ctx context.Context, id uint64) (*Tenants, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `deleted_at` is null limit 1", tenantsRows, m.table)
	var resp Tenants
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

func (m *customTenantsModel) FindOneByTenantName(ctx context.Context, tenantName string) (*Tenants, error) {
	query := fmt.Sprintf("select %s from %s where `tenant_name` = ? and `deleted_at` is null limit 1", tenantsRows, m.table)
	var resp Tenants
	err := m.conn.QueryRowCtx(ctx, &resp, query, tenantName)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customTenantsModel) Update(ctx context.Context, data *Tenants) error {
	data.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	query := fmt.Sprintf("update %s set `updated_at` = ?, `tenant_name` = ?, `email` = ? where `id` = ? and `deleted_at` is null", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.UpdatedAt, data.TenantName, data.Email, data.Id)
	return err
}

func (m *customTenantsModel) Delete(ctx context.Context, id uint64) error {
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

func (m *customTenantsModel) FindPage(ctx context.Context, tenantName string, pageIndex, pageSize int64) (int64, []*Tenants, error) {
	var (
		conds []string
		args  []any
		where string
	)
	conds = append(conds, "`deleted_at` is null")
	if tenantName != "" {
		conds = append(conds, "`tenant_name` like ?")
		args = append(args, "%"+tenantName+"%")
	}
	where = " where " + strings.Join(conds, " and ")

	var total int64
	countQuery := fmt.Sprintf("select count(*) from %s%s", m.table, where)
	if err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, err
	}

	offset := (pageIndex - 1) * pageSize
	listQuery := fmt.Sprintf("select %s from %s%s order by `updated_at` desc limit ? offset ?", tenantsRows, m.table, where)
	listArgs := append(append([]any{}, args...), pageSize, offset)

	var list []*Tenants
	if err := m.conn.QueryRowsCtx(ctx, &list, listQuery, listArgs...); err != nil {
		return 0, nil, err
	}
	return total, list, nil
}
