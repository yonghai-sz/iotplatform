package svc

import (
	"iotplatform/services/model"
	"iotplatform/services/rpc-transform/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	Model  model.PlatformModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Model:  model.NewPlatformModel(sqlx.NewMysql(c.DataSource), c.Cache),
	}
}
