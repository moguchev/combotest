package users

import (
	"combotest/internal/models"
	"context"
)

type Usecase interface {
	CreateUser(ctx context.Context, cu models.CreateUser) (models.User, error)
	GetUsers(ctx context.Context, f models.UserFilter) (uint32, []models.User, error)
	ApproveUser(ctx context.Context, id string) error
}
