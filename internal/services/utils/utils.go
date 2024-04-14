package utils

import (
	"context"

	"github.com/smakkking/avito_test/internal/httpserver/middleware"
)

func IsAdmin(ctx context.Context) bool {
	isAdmin := ctx.Value(middleware.AdminKey).(bool)
	return isAdmin
}
