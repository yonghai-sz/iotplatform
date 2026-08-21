// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package login

import (
	"context"
	"errors"
	"strings"
	"time"

	"iotplatform/pkg/utils"
	"iotplatform/services/platform-api/internal/svc"
	"iotplatform/services/platform-api/internal/types"
	"iotplatform/services/platform-api/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	errInvalidCredentials = errors.New("invalid username or password")
	errAccountDisabled    = errors.New("account is disabled")
	errRoleDisabled       = errors.New("role is disabled")
	errLoginNotAllowed    = errors.New("login is not allowed")
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	if username == "" || password == "" {
		return nil, errInvalidCredentials
	}

	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, username)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errInvalidCredentials
		}
		return nil, err
	}

	if utils.HashPassword(password, user.Salt) != user.Password {
		return nil, errInvalidCredentials
	}
	if user.Enable == "Disable" {
		return nil, errAccountDisabled
	}

	role, err := l.svcCtx.RoleModel.FindOne(l.ctx, user.RoleId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errLoginNotAllowed
		}
		return nil, err
	}
	if role.Enable == "Disable" {
		return nil, errRoleDisabled
	}

	now := time.Now().Unix()
	token, err := getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, l.svcCtx.Config.Auth.AccessExpire, jwtClaims{
		Username: username,
		UserID:   user.Id,
		RoleKey:  role.RoleKey,
		TenantID: user.TenantId,
	})
	if err != nil {
		return nil, err
	}

	if err = l.svcCtx.SessionStore.Save(l.ctx, username, token, int(l.svcCtx.Config.Auth.AccessExpire)); err != nil {
		return nil, err
	}

	return &types.LoginResp{
		Username: username,
		Token:    token,
	}, nil
}

type jwtClaims struct {
	Username string
	UserID   uint64
	RoleKey  string
	TenantID uint64
}

func getJwtToken(secretKey string, iat, seconds int64, claims jwtClaims) (string, error) {
	payload := make(jwt.MapClaims)
	payload["exp"] = iat + seconds
	payload["iat"] = iat
	payload["username"] = claims.Username
	payload["userId"] = claims.UserID
	payload["roleKey"] = claims.RoleKey
	payload["tenantId"] = claims.TenantID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(secretKey))
}
