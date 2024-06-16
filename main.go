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

const dbuser = "root"
const dbpassword = "example"
const dburi = "mongodb://localhost:27017"

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error(), "code": "500"})
	},
}

func main() {
	listenPort := flag.String("listenPort", ":5000", "The listen port of the API server")
	flag.Parse()

	credentials := options.Credential{
		Username: dbuser,
		Password: dbpassword,
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi).SetAuth(credentials))
	if err != nil {
		slog.Error("client creation", slog.String("error", err.Error()))
		return
	}
	defer client.Disconnect(context.TODO())

	app := fiber.New(config)
	apiV1 := app.Group("/api/v1")

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)

	err = app.Listen(*listenPort)
	if err != nil {
		slog.Error("server shutdown", slog.String("error", err.Error()))
		return
	}
}
