package api

import (
	"github.com/SlavaMashkov/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

var (
	user = types.User{
		FirstName: "Slava",
		LastName:  "Mashkov",
	}
)

func HandleGetUsers(c *fiber.Ctx) error {
	users := []types.User{
		user,
	}

	return c.JSON(users)
}

func HandleGetUser(c *fiber.Ctx) error {
	user.ID = c.Params("id")

	return c.JSON(user)
}
