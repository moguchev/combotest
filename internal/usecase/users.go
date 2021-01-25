package usecase

import (
	"combotest/internal/app/access"
	"combotest/internal/interfaces/users"
	"combotest/internal/models"
	"context"
	"fmt"
)

type userUsecase struct {
	repo users.Repository
}

func NewUserUscase(r users.Repository) users.Usecase {
	return &userUsecase{repo: r}
}

func (u *userUsecase) CreateUser(ctx context.Context, cu models.CreateUser) (models.User, error) {
	id, err := u.repo.CreateUser(ctx, cu)
	if err != nil {
		return models.User{}, fmt.Errorf("create user: %w", err)
	}

	cu.User.ID = id
	return cu.User, nil
}

func (u *userUsecase) GetUsers(ctx context.Context, f models.UserFilter) (uint32, []models.User, error) {
	total, err := u.repo.CountUsers(ctx, f)
	if err != nil {
		return 0, nil, fmt.Errorf("count users: %w", err)
	}
	if total == 0 {
		return 0, []models.User{}, nil
	}

	users, err := u.repo.GetUsers(ctx, f)
	if err != nil {
		return 0, nil, fmt.Errorf("get users: %w", err)
	}

	return total, users, nil
}

func (u *userUsecase) ApproveUser(ctx context.Context, id string) error {
	user, ok := access.GetUserFromCtx(ctx)
	if !ok {
		return fmt.Errorf("no user in ctx")
	}

	if user.Role != models.AdminRole {
		return fmt.Errorf("permission denied")
	}

	if err := u.repo.ApproveUser(ctx, id); err != nil {
		return fmt.Errorf("approve user: %w", err)
	}

	return nil
}
