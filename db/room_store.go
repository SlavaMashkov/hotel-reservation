package db

import (
	"context"
	"github.com/SlavaMashkov/hotel-reservation/types"
	"github.com/SlavaMashkov/hotel-reservation/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	roomCollection = "rooms"
)

type RoomStore interface {
	InsertRoom(ctx context.Context, hotelID string, room *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection
	hotelStore HotelStore
}

func NewMongoRoomStore(client *mongo.Client) *MongoRoomStore {
	collection := client.Database(NAME).Collection(roomCollection)

	return &MongoRoomStore{
		client:     client,
		collection: collection,
		hotelStore: NewMongoHotelStore(client),
	}
}

func (store *MongoRoomStore) InsertRoom(ctx context.Context, hotelID string, room *types.Room) (*types.Room, error) {
	hotelOID, err := utility.IDToMongoOID(hotelID)
	if err != nil {
		return nil, err
	}

	room.HotelID = hotelOID

	result, err := store.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.ID = result.InsertedID.(primitive.ObjectID)

	err = store.hotelStore.AppendRoom(ctx, hotelOID.Hex(), room)
	if err != nil {
		return nil, err
	}

	return room, nil
}
