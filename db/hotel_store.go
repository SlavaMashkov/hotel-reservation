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
	hotelCollection = "hotels"
)

type HotelStore interface {
	GetHotelByID(ctx context.Context, id string) (*types.Hotel, error)
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(ctx context.Context, id string, params types.UpdateHotelParams) error
}

type MongoHotelStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewHotelStoreMongo(client *mongo.Client) *MongoHotelStore {
	collection := client.Database(NAME).Collection(hotelCollection)

	return &MongoHotelStore{
		client:     client,
		collection: collection,
	}
}

func (store *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	var hotel types.Hotel

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("could not convert to ObjectID", slog.String("id", id))
		return nil, err
	}

	err = store.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel)
	if err != nil {
		return nil, err
	}

	return &hotel, nil
}

func (store *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	result, err := store.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}

	hotel.ID = result.InsertedID.(primitive.ObjectID)

	return hotel, nil
}

func (store *MongoHotelStore) UpdateHotel(ctx context.Context, id string, params types.UpdateHotelParams) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("could not convert to ObjectID", slog.String("id", id))
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
