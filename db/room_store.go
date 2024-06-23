package db

import (
	"context"
	"github.com/SlavaMashkov/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	roomCollection = "rooms"
)

type RoomStore interface {
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection
	hotelStore HotelStore
}

func NewRoomStoreMongo(client *mongo.Client) *MongoRoomStore {
	collection := client.Database(NAME).Collection(roomCollection)

	return &MongoRoomStore{
		client:     client,
		collection: collection,
		hotelStore: NewHotelStoreMongo(client),
	}
}

func (store *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	result, err := store.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.ID = result.InsertedID.(primitive.ObjectID)

	hotel, err := store.hotelStore.GetHotelByID(ctx, room.HotelID.Hex())
	if err != nil {
		return nil, err
	}

	hotel.Rooms = append(hotel.Rooms, room.ID)

	err = store.hotelStore.UpdateHotel(ctx, room.HotelID.Hex(), types.UpdateHotelParams{
		Rooms: hotel.Rooms,
	})
	if err != nil {
		return nil, err
	}

	return room, nil
}
