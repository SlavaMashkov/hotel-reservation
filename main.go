package main

import (
	"context"
	"flag"
	"github.com/SlavaMashkov/hotel-reservation/api"
	"github.com/SlavaMashkov/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{
			"error": err.Error(),
			"code":  "500",
		})
	},
}

func main() {
	listenPort := flag.String("listenPort", ":5000", "The listen port of the API server")
	flag.Parse()

	credentials := options.Credential{
		Username: db.USERNAME,
		Password: db.PASSWORD,
	}

	// TODO: extract connection initialization logic
	// Region connection initialization logic
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.URI).SetAuth(credentials))
	if err != nil {
		slog.Error("client creation", slog.String("error", err.Error()))
		return
	}
	defer client.Disconnect(context.TODO())
	// End

	// TODO: extract handlers initialization logic
	// Region handlers initialization logic

	var (
		userStore  = db.NewMongoUserStore(client)
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client)

		store = &db.Store{
			UserStore:  userStore,
			HotelStore: hotelStore,
			RoomStore:  roomStore,
		}

		userHandler  = api.NewUserHandler(store)
		hotelHandler = api.NewHotelHandler(store)

		app   = fiber.New(config)
		apiV1 = app.Group("/api/v1")
	)

	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)

	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.GetHotelRooms)
	// End

	err = app.Listen(*listenPort)
	if err != nil {
		slog.Error("server shutdown", slog.String("error", err.Error()))
		return
	}
}
