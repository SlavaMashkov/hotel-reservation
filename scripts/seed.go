package main

import (
	"context"
	"github.com/SlavaMashkov/hotel-reservation/db"
	"github.com/SlavaMashkov/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func main() {
	slog.Info("seeding the database")

	initVars()
	defer client.Disconnect(context.TODO())

	hotels := []types.Hotel{
		{
			Name:     "Hilton",
			Location: "London",
			Rooms:    make([]primitive.ObjectID, 0),
			Rating:   4,
		},
		{
			Name:     "Radisson",
			Location: "Paris",
			Rooms:    make([]primitive.ObjectID, 0),
			Rating:   5,
		},
		{
			Name:     "Inn",
			Location: "New York",
			Rooms:    make([]primitive.ObjectID, 0),
			Rating:   3,
		},
	}

	for _, hotel := range hotels {
		_, err := seedHotel(&hotel)
		if err != nil {
			return
		}
	}

	slog.Info("database seeded")
}

func initVars() {
	credentials := options.Credential{
		Username: db.USERNAME,
		Password: db.PASSWORD,
	}

	ctx = context.Background()
	var err error

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(db.URI).SetAuth(credentials))
	if err != nil {
		slog.Error("client creation", slog.String("error", err.Error()))
		return
	}

	if err = client.Database(db.NAME).Drop(ctx); err != nil {
		slog.Error("database dropping", slog.String("error", err.Error()))
		return
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client)
}

func seedHotel(hotel *types.Hotel) (*types.Hotel, error) {
	insertedHotel, err := hotelStore.InsertHotel(ctx, hotel)
	if err != nil {
		slog.Error("hotel insertion", slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info("hotel inserted", slog.String("id", insertedHotel.ID.Hex()))

	rooms := []types.Room{
		{
			Type:      types.Deluxe,
			BasePrice: 100,
		},
		{
			Type:      types.Single,
			BasePrice: 50,
		},
	}

	err = seedRooms(rooms, insertedHotel)
	if err != nil {
		return nil, err
	}

	return insertedHotel, nil
}

func seedRooms(rooms []types.Room, hotel *types.Hotel) error {
	for _, room := range rooms {
		_, err := roomStore.InsertRoom(ctx, hotel.ID.Hex(), &room)
		if err != nil {
			slog.Error("room insertion", slog.String("error", err.Error()))
			return err
		}

		slog.Info("room inserted", slog.String("id", room.ID.Hex()))
	}

	return nil
}
