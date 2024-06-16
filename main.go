package main

import (
	"context"
	"flag"
	"github.com/SlavaMashkov/hotel-reservation/api"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

const dbuser = "root"
const dbpassword = "example"
const dburi = "mongodb://localhost:27017"

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

	app := fiber.New()

	apiV1 := app.Group("/api/v1")

	apiV1.Get("/user", api.HandleGetUsers)
	apiV1.Get("/user/:id", api.HandleGetUser)

	err := app.Listen(*listenPort)
	if err != nil {
		slog.Error("server shutdown", slog.String("error", err.Error()))
		return
	}
}
