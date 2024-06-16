package main

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func main() {
	app := fiber.New()
	slog.Info("start server", slog.Int("port", 5000))
	app.Get("/ping", handlePing())
	app.Listen(":5000")
}

func handlePing() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{"message": "pong"})
	}
}
