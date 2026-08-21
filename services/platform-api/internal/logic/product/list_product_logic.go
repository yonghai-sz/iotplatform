// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package product

import (
	"context"

	"iotplatform/services/platform-api/internal/svc"
	"iotplatform/services/platform-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductLogic {
	return &ListProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductLogic) ListProduct(req *types.ListProductReq) (resp *types.ListProductResp, err error) {
	pageIndex := req.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	total, list, err := l.svcCtx.ProductModel.FindPage(l.ctx, req.ProductName, pageIndex, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]types.ProductInfo, 0, len(list))
	for _, item := range list {
		items = append(items, toProductInfo(item))
	}

	return &types.ListProductResp{
		Total: total,
		List:  items,
	}, nil
}
