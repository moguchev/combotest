package events

import (
	"combotest/internal/models"
	"context"
)

// Usecase - usecase lvl for events
type Usecase interface {
	UpdateEvent(ctx context.Context, id string, fields models.UpdateEventFields) error
	GetEvents(ctx context.Context, filter models.EventsFilter) (uint32, []models.Event, error)
	SetIncedent(ctx context.Context, ids []string) error
}
