package middleware

import (
	"context"
	"net/http"
	"redditclone/pkg/session"
)

type ctxKey int

const SessionKey ctxKey = 1

func AuthMiddleware(sm *session.SessionJWT, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := sm.Check(r)
		if err != nil {
			http.Error(w, "No auth", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), SessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
