package api

import (
	"github.com/SlavaMashkov/hotel-reservation/db"
	"github.com/SlavaMashkov/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStoreGetter db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStoreGetter,
	}
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := types.NewUser(&params)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(insertedUser)
}
