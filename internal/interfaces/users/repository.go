package users

import (
	"combotest/internal/models"
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, cu models.CreateUser) (string, error)
	CountUsers(ctx context.Context, f models.UserFilter) (uint32, error)
	GetUsers(ctx context.Context, f models.UserFilter) ([]models.User, error)
	ApproveUser(ctx context.Context, id string) error
	GetUserByAuthInfo(ctx context.Context, a models.AuthInfo) (models.User, error)
}
