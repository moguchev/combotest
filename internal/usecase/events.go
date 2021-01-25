package usecase

import (
	"combotest/internal/app/access"
	"combotest/internal/app/closessl"
	"combotest/internal/interfaces/events"
	"combotest/internal/models"
	"context"
	"fmt"
)

type evensUsecase struct {
	repo events.Repository
}

func NewEventUscase(r events.Repository) events.Usecase {
	return &evensUsecase{repo: r}
}

func (u *evensUsecase) UpdateEvent(ctx context.Context, id string, fields models.UpdateEventFields) error {
	return nil
}

func (u *evensUsecase) GetEvents(ctx context.Context, filter models.EventsFilter) (uint32, []models.Event, error) {
	user, ok := access.GetUserFromCtx(ctx)
	if !ok {
		return 0, nil, fmt.Errorf("no user in ctx")
	}

	total, err := u.repo.CountEvents(ctx, filter)
	if err != nil {
		return 0, nil, fmt.Errorf("count events: %w", err)
	}

	events, err := u.repo.GetEvents(ctx, filter)
	if err != nil {
		return 0, nil, fmt.Errorf("get events: %w", err)
	}

	if user.Role == models.AnalystRole && user.Confirmed {
		for i := range events {
			events[i].Message = string(closessl.Decrypt([]byte(events[i].Message)))
		}
	}

	return total, events, nil
}
