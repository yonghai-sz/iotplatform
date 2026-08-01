// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package product

import (
	"context"
	"errors"

	"iotplatform/services/api/internal/svc"
	"iotplatform/services/api/internal/types"
	"iotplatform/services/api/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductLogic) GetProduct(req *types.ProductIdPathReq) (resp *types.ProductInfo, err error) {
	entity, err := l.svcCtx.ProductModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	info := toProductInfo(entity)
	return &info, nil
}
