package usecase

import (
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
	return 0, nil, nil
}

func (u *userUsecase) ApproveUser(ctx context.Context, id string) error {
	return nil
}
