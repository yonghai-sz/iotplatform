package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"iotplatform/services/api/internal/session"
)

func TestSessionMiddleware_Handle(t *testing.T) {
	store := session.NewStore(redis.New(miniredis.RunT(t).Addr()))
	handler := NewSessionMiddleware(store).Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("missing username", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("session missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		req = req.WithContext(context.WithValue(req.Context(), "username", "alice"))
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("token mismatch", func(t *testing.T) {
		assert.NoError(t, store.Save(context.Background(), "alice", "token-1", 3600))
		req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		req.Header.Set("Authorization", "Bearer other-token")
		req = req.WithContext(context.WithValue(req.Context(), "username", "alice"))
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("success", func(t *testing.T) {
		assert.NoError(t, store.Save(context.Background(), "alice", "token-1", 3600))
		req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		req = req.WithContext(context.WithValue(req.Context(), "username", "alice"))
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
