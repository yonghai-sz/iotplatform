// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"errors"
	"strings"

	"iotplatform/pkg/utils"
	"iotplatform/services/platform-api/internal/svc"
	"iotplatform/services/platform-api/internal/types"
	"iotplatform/services/platform-api/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserPasswordLogic {
	return &UpdateUserPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserPasswordLogic) UpdateUserPassword(req *types.UpdateUserPasswordReq) error {
	password := strings.TrimSpace(req.Password)
	if password == "" {
		return errors.New("password is required")
	}

	entity, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	salt := utils.GenerateCryptoRandString(6)
	hashed := utils.HashPassword(password, salt)

	entity.Password = hashed
	entity.Salt = salt

	return l.svcCtx.UserModel.Update(l.ctx, entity)
}
