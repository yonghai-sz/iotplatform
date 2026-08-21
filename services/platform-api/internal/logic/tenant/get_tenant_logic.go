// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package tenant

import (
	"context"
	"errors"

	"iotplatform/services/platform-api/internal/svc"
	"iotplatform/services/platform-api/internal/types"
	"iotplatform/services/platform-api/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTenantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTenantLogic {
	return &GetTenantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTenantLogic) GetTenant(req *types.TenantIdPathReq) (resp *types.TenantInfo, err error) {
	entity, err := l.svcCtx.TenantsModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}

	info := toTenantInfo(entity)
	return &info, nil
}
