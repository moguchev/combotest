package access

import (
	"combotest/internal/models"
	"context"
)

type ctxKey = struct {
	key string
}

var roleCtxKey = ctxKey{"role"}

func GetRoleFromCtx(ctx context.Context) (models.Role, bool) {
	role, ok := ctx.Value(roleCtxKey).(models.Role)
	if !ok {
		return models.Anonymous, false
	}

	return role, true
}

func SetRoleToCtx(ctx context.Context, role models.Role) context.Context {
	return context.WithValue(ctx, roleCtxKey, role)
}
