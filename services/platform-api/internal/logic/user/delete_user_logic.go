// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"errors"

	"iotplatform/services/platform-api/internal/svc"
	"iotplatform/services/platform-api/internal/types"
	"iotplatform/services/platform-api/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.UserUsernamePathReq) error {
	entity, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if err = l.svcCtx.UserModel.Delete(l.ctx, entity.Id); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	return nil
}
