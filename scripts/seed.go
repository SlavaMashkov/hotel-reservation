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

func main() {
	slog.Info("seeding the database")

	credentials := options.Credential{
		Username: db.USERNAME,
		Password: db.PASSWORD,
	}

	ctx := context.Background()

	// TODO: extract connection initialization logic
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.URI).SetAuth(credentials))
	if err != nil {
		slog.Error("client creation", slog.String("error", err.Error()))
		return
	}
	defer client.Disconnect(context.TODO())

	hotelStore := db.NewHotelStoreMongo(client)
	roomStore := db.NewRoomStoreMongo(client)

	hotel := types.Hotel{
		Name:     "Hilton",
		Location: "London",
		Rooms:    make([]primitive.ObjectID, 0),
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		slog.Error("hotel insertion", slog.String("error", err.Error()))
		return
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

	for _, room := range rooms {
		_, err = roomStore.InsertRoom(ctx, insertedHotel.ID.Hex(), &room)
		if err != nil {
			slog.Error("room insertion", slog.String("error", err.Error()))
			return
		}

		slog.Info("room inserted", slog.String("id", room.ID.Hex()))
	}

	slog.Info("database seeded")
}
