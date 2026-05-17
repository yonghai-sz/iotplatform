package svc

import (
	"iotplatform/services/rpc-transform/internal/config"
	"iotplatform/services/rpc-transform/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	Model  model.ExampleModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Model:  model.NewExampleModel(sqlx.NewMysql(c.DataSource), c.Cache),
	}
}
