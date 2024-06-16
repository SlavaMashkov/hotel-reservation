package main

import (
	"flag"
	"github.com/SlavaMashkov/hotel-reservation/api"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func main() {
	listenPort := flag.String("listenPort", ":5000", "The listen port of the API server")
	flag.Parse()

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
