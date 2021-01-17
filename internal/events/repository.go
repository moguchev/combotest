package events

import (
	"combotest/internal/models"
	"context"
)

// Repository - repo lvl for events
type Repository interface {
	InsertEvent(ctx context.Context, e models.Event) error
	InsertEvents(ctx context.Context, es []models.Event) error
	UpdateEvent(ctx context.Context, id string, fields models.UpdateEventFields) error
	GetEvents(ctx context.Context, filter models.EventsFilter) ([]models.Event, error)
}
