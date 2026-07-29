package httpadapter

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	contextKeyUserID contextKey = "userID"
	contextKeyEmail  contextKey = "email"
)

func (r *Router) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, http.StatusUnauthorized, "인증 토큰이 필요합니다.")
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := r.tokens.ParseAccessToken(tokenStr)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "유효하지 않은 토큰입니다.")
			return
		}
		ctx := context.WithValue(req.Context(), contextKeyUserID, claims["sub"])
		ctx = context.WithValue(ctx, contextKeyEmail, claims["email"])
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
