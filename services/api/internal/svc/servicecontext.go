// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"iotplatform/services/api/internal/config"
	"iotplatform/services/api/model"
	"iotplatform/services/rpc-transform/transformer"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	Transformer  transformer.Transformer
	TenantsModel model.TenantsModel
	ProductModel model.ProductModel
	RoleModel    model.RoleModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:       c,
		Transformer:  transformer.NewTransformer(zrpc.MustNewClient(c.Transform)),
		TenantsModel: model.NewTenantsModel(conn),
		ProductModel: model.NewProductModel(conn),
		RoleModel:    model.NewRoleModel(conn),
	}
}
