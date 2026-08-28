package login

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"iot-zero/pkg/utils"
	"iot-zero/services/platform-api/internal/config"
	"iot-zero/services/platform-api/internal/session"
	"iot-zero/services/platform-api/internal/svc"
	"iot-zero/services/platform-api/internal/types"
	"iot-zero/services/platform-api/model"
)

func newLoginLogicWithMock(t *testing.T) (*LoginLogic, sqlmock.Sqlmock, *session.Store) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	conn := sqlx.NewSqlConnFromDB(db)
	cfg := config.Config{}
	cfg.Auth.AccessSecret = "test-secret"
	cfg.Auth.AccessExpire = 3600
	store := session.NewStore(redis.New(miniredis.RunT(t).Addr()))

	return NewLoginLogic(context.Background(), &svc.ServiceContext{
		Config:       cfg,
		UserModel:    model.NewUserModel(conn),
		RoleModel:    model.NewRoleModel(conn),
		SessionStore: store,
	}), mock, store
}

func TestLoginLogic_Login(t *testing.T) {
	const (
		username = "alice"
		password = "secret"
		salt     = "salt"
		roleType = "admin"
	)
	hashed := utils.HashPassword(password, salt)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		ast := assert.New(t)
		l, mock, store := newLoginLogicWithMock(t)

		mock.ExpectQuery("select .+ from `user` where `username` = \\? and `deleted_at` is null").
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "password", "salt", "enable", "role_id", "tenant_id"}).
				AddRow(uint64(8), now, now, nil, username, hashed, salt, "Enable", uint64(2), uint64(1)))
		mock.ExpectQuery("select .+ from `role` where `id` = \\? and `deleted_at` is null").
			WithArgs(uint64(2)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "role_type", "role_name", "enable", "tenant_id"}).
				AddRow(uint64(2), now, now, nil, roleType, "Admin", "Enable", uint64(1)))

		resp, err := l.Login(&types.LoginReq{Username: username, Password: password})
		ast.NoError(err)
		ast.Equal(username, resp.Username)
		ast.NotEmpty(resp.Token)

		parsed, err := jwt.Parse(resp.Token, func(token *jwt.Token) (any, error) {
			return []byte("test-secret"), nil
		})
		ast.NoError(err)
		claims := parsed.Claims.(jwt.MapClaims)
		ast.Equal(username, claims["username"])
		ast.Equal(roleType, claims["roleType"])
		ast.Equal(float64(8), claims["userId"])
		ast.Equal(float64(1), claims["tenantId"])

		matched, err := store.Match(context.Background(), username, resp.Token)
		ast.NoError(err)
		ast.True(matched)
	})

	t.Run("wrong password", func(t *testing.T) {
		ast := assert.New(t)
		l, mock, _ := newLoginLogicWithMock(t)

		mock.ExpectQuery("select .+ from `user` where `username` = \\? and `deleted_at` is null").
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "password", "salt", "enable", "role_id", "tenant_id"}).
				AddRow(uint64(8), now, now, nil, username, hashed, salt, "Enable", uint64(2), uint64(1)))

		_, err := l.Login(&types.LoginReq{Username: username, Password: "wrong"})
		ast.ErrorIs(err, errInvalidCredentials)
	})

	t.Run("user not found", func(t *testing.T) {
		ast := assert.New(t)
		l, mock, _ := newLoginLogicWithMock(t)

		mock.ExpectQuery("select .+ from `user` where `username` = \\? and `deleted_at` is null").
			WithArgs(username).
			WillReturnError(sqlx.ErrNotFound)

		_, err := l.Login(&types.LoginReq{Username: username, Password: password})
		ast.ErrorIs(err, errInvalidCredentials)
	})

	t.Run("account disabled", func(t *testing.T) {
		ast := assert.New(t)
		l, mock, _ := newLoginLogicWithMock(t)

		mock.ExpectQuery("select .+ from `user` where `username` = \\? and `deleted_at` is null").
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "password", "salt", "enable", "role_id", "tenant_id"}).
				AddRow(uint64(8), now, now, nil, username, hashed, salt, "Disable", uint64(2), uint64(1)))

		_, err := l.Login(&types.LoginReq{Username: username, Password: password})
		ast.ErrorIs(err, errAccountDisabled)
	})

	t.Run("role disabled", func(t *testing.T) {
		ast := assert.New(t)
		l, mock, _ := newLoginLogicWithMock(t)

		mock.ExpectQuery("select .+ from `user` where `username` = \\? and `deleted_at` is null").
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "password", "salt", "enable", "role_id", "tenant_id"}).
				AddRow(uint64(8), now, now, nil, username, hashed, salt, "Enable", uint64(2), uint64(1)))
		mock.ExpectQuery("select .+ from `role` where `id` = \\? and `deleted_at` is null").
			WithArgs(uint64(2)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "role_type", "role_name", "enable", "tenant_id"}).
				AddRow(uint64(2), now, now, nil, roleType, "Admin", "Disable", uint64(1)))

		_, err := l.Login(&types.LoginReq{Username: username, Password: password})
		ast.ErrorIs(err, errRoleDisabled)
	})
}
