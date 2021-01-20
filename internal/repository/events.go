package repository

import (
	"combotest/internal/interfaces/events"
	"combotest/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type eventsRepository struct {
	db *mongo.Database
}

const (
	eventCollection = "events"
)

// NewEventsRepository - return implimentation of events.Repository interface
func NewEventsRepository(db *mongo.Database) events.Repository {
	return &eventsRepository{db: db}
}

func (r *eventsRepository) InsertEvent(ctx context.Context, e models.Event) error {
	collection := r.db.Collection(eventCollection)

	_, err := collection.InsertOne(ctx, e)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}

	return nil
}

func (r *eventsRepository) InsertEvents(ctx context.Context, es []models.Event) error {
	collection := r.db.Collection(eventCollection)

	docs := make([]interface{}, 0, len(es))
	for i := range es {
		es[i].ID = primitive.NewObjectID().String()
		docs = append(docs, es[i])
	}

	res, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("insert events: %w", err)
	}

	fmt.Println(res.InsertedIDs)

	return nil
}

func (r *eventsRepository) UpdateEvent(ctx context.Context, id string, fields models.UpdateEventFields) error {
	collection := r.db.Collection(eventCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("object id from hex: %w", err)
	}

	update := bson.D{}

	if fields.IsIncident != nil {
		update = append(update, bson.E{Key: "is_incident", Value: *fields.IsIncident})
	}

	if len(update) == 0 {
		return nil
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}}
	_, err = collection.UpdateOne(ctx, filter,
		bson.M{"$set": update},
	)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}

	return nil
}

func getEventOptions(filter models.EventsFilter) *options.FindOptions {
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	if filter.Limit != nil {
		opts.SetLimit(int64(*filter.Limit))
	}

	if filter.Offset != nil {
		opts.SetSkip(int64(*filter.Offset))
	}

	return opts
}

func getEventFilter(filter models.EventsFilter) bson.M {
	f := bson.M{}

	if filter.Incident != nil {
		f["is_incident"] = *filter.Incident
	}

	if filter.SystemName != nil {
		f["system_name"] = *filter.SystemName
	}

	if filter.EventID != nil {
		f["event_id"] = *filter.EventID
	}

	if filter.CreatedAt != nil {
		m := bson.M{}
		if filter.CreatedAt.Begin != nil {
			m["$gte"] = *filter.CreatedAt.Begin
		}

		if filter.CreatedAt.End != nil {
			m["$lt"] = *filter.CreatedAt.End
		}
	}

	return f
}

func (r *eventsRepository) GetEvents(ctx context.Context, f models.EventsFilter) ([]models.Event, error) {
	collection := r.db.Collection(eventCollection)

	options := getEventOptions(f)
	filter := getEventFilter(f)

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	defer cursor.Close(ctx)

	events := make([]models.Event, 0, 2)

	for cursor.Next(ctx) {
		var event models.Event
		if err = cursor.Decode(&event); err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}
		events = append(events, event)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return events, nil
}
