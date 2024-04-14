package middleware

import (
	"context"
	"net/http"
)

type ctxKeyAdmin string

var AdminKey ctxKeyAdmin = "admin"

func Authorization(next http.Handler) http.Handler {
	// другой вопрос - как мы здесь должны проверять авторизацию - может быть нужен отдельный сервис или достаточно захардкодить токены?

	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token := r.Header.Get("token")

		switch token {
		case "admin":
			ctx = context.WithValue(ctx, AdminKey, true)
		case "user":
			ctx = context.WithValue(ctx, AdminKey, false)
		default:
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
