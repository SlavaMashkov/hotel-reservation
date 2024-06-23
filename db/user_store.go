package db

import (
	"context"
	"github.com/SlavaMashkov/hotel-reservation/types"
	"github.com/SlavaMashkov/hotel-reservation/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollection = "users"
)

type UserStore interface {
	GetUsers(ctx context.Context) ([]*types.User, error)
	GetUser(ctx context.Context, userID string) (*types.User, error)
	InsertUser(ctx context.Context, user *types.User) (*types.User, error)
	UpdateUser(ctx context.Context, userID string, params types.UpdateUserParams) error
	DeleteUser(ctx context.Context, userID string) error
}

type MongoUserStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	collection := client.Database(NAME).Collection(userCollection)

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

func (store *MongoUserStore) GetUser(ctx context.Context, userID string) (*types.User, error) {
	var user types.User

	oid, err := utility.IDToMongoOID(userID)
	if err != nil {
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

func (store *MongoUserStore) UpdateUser(ctx context.Context, userID string, params types.UpdateUserParams) error {
	oid, err := utility.IDToMongoOID(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	update := params.ToBSON()

	result, err := store.collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (store *MongoUserStore) DeleteUser(ctx context.Context, userID string) error {
	userOID, err := utility.IDToMongoOID(userID)
	if err != nil {
		return err
	}

	result, err := store.collection.DeleteOne(ctx, bson.M{"_id": userOID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
