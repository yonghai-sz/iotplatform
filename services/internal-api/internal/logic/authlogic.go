// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"time"

	"iotplatform/services/internal-api/internal/svc"
	"iotplatform/services/internal-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthLogic {
	return &AuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Auth implements EMQX HTTP password-based authentication (allow-all stub).
// See https://docs.emqx.com/en/emqx/v5.8/access-control/authn/http.html
func (l *AuthLogic) Auth(req *types.AuthReq) (*types.AuthResp, error) {
	l.Infow("emqx authn", logx.Field("clientid", req.Clientid), logx.Field("username", req.Username))
	return &types.AuthResp{
		Result:      "allow",
		IsSuperuser: false,
		ExpireAt:    time.Now().Add(1 * time.Hour).Unix(),
		// Acl: []types.AuthAclRule{
		// 	{
		// 		Permission: "publish",
		// 		Action:     "allow",
		// 		Topic:      "#",
		// 	},
		// 	{
		// 		Permission: "subscribe",
		// 		Action:     "allow",
		// 		Topic:      "#",
		// 	},
		// },
		// ClientAttrs: map[string]string{
		// 	"username": req.Username,
		// 	"password": req.Password,
		// },
	}, nil
}
