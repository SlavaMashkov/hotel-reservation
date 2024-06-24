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
	GetHotels(ctx context.Context) ([]*types.Hotel, error)
	GetHotel(ctx context.Context, hotelID string) (*types.Hotel, error)
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(ctx context.Context, hotelID string, params types.UpdateHotelParams) error
	AppendRoom(ctx context.Context, hotelID string, room *types.Room) error
}

type MongoHotelStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	collection := client.Database(NAME).Collection(hotelCollection)

	return &MongoHotelStore{
		client:     client,
		collection: collection,
	}
}

func (store *MongoHotelStore) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	cur, err := store.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel
	if err = cur.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (store *MongoHotelStore) GetHotel(ctx context.Context, hotelID string) (*types.Hotel, error) {
	var hotel types.Hotel

	oid, err := utility.IDToMongoOID(hotelID)
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

func (store *MongoHotelStore) UpdateHotel(ctx context.Context, hotelID string, params types.UpdateHotelParams) error {
	hotelOID, err := utility.IDToMongoOID(hotelID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": hotelOID}

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

func (store *MongoHotelStore) AppendRoom(ctx context.Context, hotelID string, room *types.Room) error {
	hotelOID, err := utility.IDToMongoOID(hotelID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": hotelOID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}

	_, err = store.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
