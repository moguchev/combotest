package repository

import (
	"combotest/internal/interfaces/users"
	"combotest/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	usersCollection = "users"
)

type usersRepository struct {
	db *mongo.Database
}

func NewUsersRepository(db *mongo.Database) users.Repository {
	return &usersRepository{db: db}
}

func (r *usersRepository) CreateUser(ctx context.Context, cu models.CreateUser) (string, error) {
	collection := r.db.Collection(usersCollection)

	id := primitive.NewObjectID()

	if cu.ID != "" {
		old, err := primitive.ObjectIDFromHex(cu.ID)
		if err != nil {
			id = old
		}
	}

	u := struct {
		ID        primitive.ObjectID `json:"id"        bson:"_id"`
		Role      models.Role        `json:"role"      bson:"role"`
		Confirmed bool               `json:"confirmed" bson:"confirmed"`
		Login     string             `json:"login"     bson:"login"`
		Password  string             `json:"password"  bson:"password"`
	}{
		ID:        id,
		Role:      cu.Role,
		Confirmed: cu.Confirmed,
		Login:     cu.Login,
		Password:  cu.Password,
	}

	doc := interface{}(u)

	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("insert events: %w", err)
	}

	return id.Hex(), nil
}

func (r *usersRepository) ApproveUser(ctx context.Context, id string) error {
	collection := r.db.Collection(usersCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("object id from hex: %w", err)
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}}
	update := bson.D{{Key: "confirmed", Value: true}}

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("approve user: %w", err)
	}

	return nil
}

func (r *usersRepository) CountUsers(ctx context.Context, f models.UserFilter) (uint32, error) {
	collection := r.db.Collection(usersCollection)

	search := getUserFilter(f)

	count, err := collection.CountDocuments(ctx, search)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}

	return uint32(count), nil
}

func getUserOptions(filter models.UserFilter) *options.FindOptions {
	opts := options.Find()

	if filter.Limit != nil {
		opts = opts.SetLimit(int64(*filter.Limit))
	}

	if filter.Offset != nil {
		opts = opts.SetSkip(int64(*filter.Offset))
	}

	return opts
}

func getUserFilter(filter models.UserFilter) bson.M {
	f := bson.M{}

	if filter.Confirmed != nil {
		f["confirmed"] = *filter.Confirmed
	}

	if filter.Role != nil {
		f["role"] = *filter.Role
	}

	return f
}

func (r *usersRepository) GetUsers(ctx context.Context, f models.UserFilter) ([]models.User, error) {
	collection := r.db.Collection(usersCollection)

	options := getUserOptions(f)
	filter := getUserFilter(f)

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	defer cursor.Close(ctx)

	users := make([]models.User, 0, 2)

	for cursor.Next(ctx) {
		var usr models.User
		if err = cursor.Decode(&usr); err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}
		users = append(users, usr)
	}

	return users, nil
}

func (r *usersRepository) GetUserByAuthInfo(ctx context.Context, a models.AuthInfo) (models.User, error) {
	collection := r.db.Collection(usersCollection)

	filter := bson.M{"login": a.Login, "password": a.Password}

	fmt.Println(filter, usersCollection, a)

	var user models.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("find user: %w", err)
	}

	return user, nil
}
