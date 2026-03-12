package middleware

import (
	"context"
	"net/http"
	"strings"

	"inventari/api/internal/session"
)

type ctxKey string

// CtxUser is the context key used to store the authenticated user.
const CtxUser ctxKey = "user"

func Auth(store *session.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && r.URL.Path == "/auth/login" {
				next.ServeHTTP(w, r)
				return
			}

			bearer := r.Header.Get("Authorization")
			token := strings.TrimPrefix(bearer, "Bearer ")
			if token == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			user, ok := store.Get(token)
			if !ok {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CtxUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
