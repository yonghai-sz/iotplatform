package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ExampleModel = (*customExampleModel)(nil)

type (
	// ExampleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customExampleModel.
	ExampleModel interface {
		exampleModel
	}

	customExampleModel struct {
		*defaultExampleModel
	}
)

// NewExampleModel returns a model for the database table.
func NewExampleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ExampleModel {
	return &customExampleModel{
		defaultExampleModel: newExampleModel(conn, c, opts...),
	}
}
