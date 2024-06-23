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
	hotelCollection = "hotels"
)

type HotelStore interface {
	GetHotelByID(ctx context.Context, id string) (*types.Hotel, error)
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(ctx context.Context, id string, params types.UpdateHotelParams) error
	AddRoom(ctx context.Context, id string, room *types.Room) error
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

	oid, err := utility.IDToMongoOID(id)
	if err != nil {
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
	oid, err := utility.IDToMongoOID(id)
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

func (store *MongoHotelStore) AddRoom(ctx context.Context, id string, room *types.Room) error {
	oid, err := utility.IDToMongoOID(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}

	_, err = store.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
