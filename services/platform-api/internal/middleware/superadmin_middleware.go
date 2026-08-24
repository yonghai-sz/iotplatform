// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package middleware

import (
	"net/http"

	"iot-zero/services/platform-api/internal/session"
)

type SuperAdminMiddleware struct{}

func NewSuperAdminMiddleware() *SuperAdminMiddleware {
	return &SuperAdminMiddleware{}
}

func (m *SuperAdminMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleKey, ok := session.RoleKeyFromContext(r.Context())
		if !ok || roleKey != session.RoleKeySuperAdmin {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
