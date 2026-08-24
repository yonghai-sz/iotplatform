package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"iot-zero/services/platform-api/internal/session"
)

func TestSuperAdminMiddleware_Handle(t *testing.T) {
	handler := NewSuperAdminMiddleware().Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("missing role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/products", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("non super admin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/products", nil)
		req = req.WithContext(context.WithValue(req.Context(), "roleKey", "admin"))
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("super admin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/products", nil)
		req = req.WithContext(context.WithValue(req.Context(), "roleKey", session.RoleKeySuperAdmin))
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
