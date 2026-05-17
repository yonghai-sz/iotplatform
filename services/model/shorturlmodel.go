package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PlatformModel = (*customPlatformModel)(nil)

type (
	// PlatformModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPlatformModel.
	PlatformModel interface {
		platformModel
	}

	customPlatformModel struct {
		*defaultPlatformModel
	}
)

// NewPlatformModel returns a model for the database table.
func NewPlatformModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PlatformModel {
	return &customPlatformModel{
		defaultPlatformModel: newPlatformModel(conn, c, opts...),
	}
}
