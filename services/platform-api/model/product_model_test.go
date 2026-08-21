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

func newMockProductModel(t *testing.T) (ProductModel, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return NewProductModel(sqlx.NewSqlConnFromDB(db)), mock
}

func TestProductModel_InsertAndFindOne(t *testing.T) {
	ast := assert.New(t)
	m, mock := newMockProductModel(t)
	ctx := context.Background()

	mock.ExpectExec("insert into `product`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "P001", "Sensor").
		WillReturnResult(sqlmock.NewResult(5, 1))

	result, err := m.Insert(ctx, &Product{
		ProductCode: "P001",
		ProductName: "Sensor",
	})
	ast.NoError(err)
	id, err := result.LastInsertId()
	ast.NoError(err)
	ast.Equal(int64(5), id)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "product_code", "product_name"}).
		AddRow(uint64(5), now, now, nil, "P001", "Sensor")
	mock.ExpectQuery("select .+ from `product` where `id` = \\? and `deleted_at` is null").
		WithArgs(uint64(5)).
		WillReturnRows(rows)

	got, err := m.FindOne(ctx, 5)
	ast.NoError(err)
	ast.Equal("P001", got.ProductCode)
	ast.Equal("Sensor", got.ProductName)
}

func TestProductModel_FindPage(t *testing.T) {
	ast := assert.New(t)
	m, mock := newMockProductModel(t)
	ctx := context.Background()

	mock.ExpectQuery("select count\\(\\*\\) from `product` where `deleted_at` is null and `product_name` like \\?").
		WithArgs("%Sensor%").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))

	now := time.Now()
	mock.ExpectQuery("select .+ from `product` where `deleted_at` is null and `product_name` like \\? order by `updated_at` desc limit \\? offset \\?").
		WithArgs("%Sensor%", int64(20), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "product_code", "product_name"}).
			AddRow(uint64(1), now, now, nil, "P001", "Sensor"))

	total, list, err := m.FindPage(ctx, "Sensor", 1, 20)
	ast.NoError(err)
	ast.Equal(int64(1), total)
	ast.Len(list, 1)
	ast.Equal("Sensor", list[0].ProductName)
}

func TestProductModel_SoftDelete(t *testing.T) {
	ast := assert.New(t)
	m, mock := newMockProductModel(t)
	ctx := context.Background()

	mock.ExpectExec(fmt.Sprintf("update %s set `deleted_at`", "`product`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), uint64(3)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	ast.NoError(m.Delete(ctx, 3))

	mock.ExpectExec(fmt.Sprintf("update %s set `deleted_at`", "`product`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), uint64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	ast.ErrorIs(m.Delete(ctx, 99), ErrNotFound)
}
