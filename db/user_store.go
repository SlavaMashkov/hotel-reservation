package db

import (
	"context"
	"github.com/SlavaMashkov/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

const (
	userCollection = "users"
)

type UserStore interface {
	GetUserByID(ctx context.Context, id string) (*types.User, error)
	GetUsers(ctx context.Context) ([]types.User, error)
}

type MongoUserStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	collection := client.Database(dbname).Collection(userCollection)

	return &MongoUserStore{
		client:     client,
		collection: collection,
	}
}

func (store *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	var user types.User

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("could not convert to ObjectID", slog.String("id", id))
		return nil, err
	}

	err = store.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (store *MongoUserStore) GetUsers(ctx context.Context) ([]types.User, error) {
	//TODO implement me
	panic("implement me")
}