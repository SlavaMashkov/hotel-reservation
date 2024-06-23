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
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	UpdateUser(ctx context.Context, id string, params types.UpdateUserParams) error
	DeleteUser(context.Context, string) error
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

func (store *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := store.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var users []*types.User
	if err = cur.All(ctx, &users); err != nil {
		return []*types.User{}, nil
	}

	return users, nil
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

func (store *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	result, err := store.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (store *MongoUserStore) UpdateUser(ctx context.Context, id string, params types.UpdateUserParams) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("could not convert to ObjectID", slog.String("id", id))
		return err
	}

	filter := bson.M{"_id": oid}

	update := bson.M{}

	if params.FirstName != "" {
		update["firstName"] = params.FirstName
	}
	if params.LastName != "" {
		update["lastName"] = params.LastName
	}

	result, err := store.collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (store *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("could not convert to ObjectID", slog.String("id", id))
		return err
	}

	result, err := store.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
