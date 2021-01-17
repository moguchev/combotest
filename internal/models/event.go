package models

import (
	"time"
)

type Event struct {
	ID         string    `json:"id"                    bson:"_id"`
	EventID    int       `json:"EventID"               bson:"event_id"`
	Created    time.Time `json:"Created"               bson:"created_at"`
	SystemName string    `json:"SystemName"            bson:"system_name"`
	Message    string    `json:"Message"               bson:"message"`
	IsIncident *bool     `json:"is_incident,omitempty" bson:"is_incident,omitempty"`
}

type UpdateEventFields struct {
	IsIncident *bool `json:"message,omitempty"`
}

// TimeInterval - BETWEEN [Begin, End]
type TimeInterval struct {
	Begin *time.Time
	End   *time.Time
}

type EventsFilter struct {
	Limit  *uint32
	Offset *uint32

	EventID    *int    // EQUAL
	SystemName *string // EQUAL
	Incident   *bool   // EQUAL

	CreatedAt *TimeInterval // BETWEEN
}
