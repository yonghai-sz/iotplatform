// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package middleware

import (
	"net/http"

	"iotplatform/services/api/internal/session"
)

type SessionMiddleware struct {
	store *session.Store
}

func NewSessionMiddleware(store *session.Store) *SessionMiddleware {
	return &SessionMiddleware{store: store}
}

func (m *SessionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := session.UsernameFromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := session.BearerToken(r)
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ok, err := m.store.Match(r.Context(), username, token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
