package api

import (
	"errors"
	"github.com/SlavaMashkov/hotel-reservation/db"
	"github.com/SlavaMashkov/hotel-reservation/types"
	"github.com/SlavaMashkov/hotel-reservation/utility"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var queryParams types.QueryHotelParams
	if err := c.QueryParser(&queryParams); err != nil {
		return err
	}

	slog.Info("query params", slog.Any("value", queryParams))

	hotels, err := h.store.HotelStore.GetHotels(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")

	hotel, err := h.store.HotelStore.GetHotel(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "hotel not found",
				"id":      id,
			})
		}

		return err
	}

	return c.JSON(hotel)
}

func (h *HotelHandler) GetHotelRooms(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := utility.IDToMongoOID(id)
	if err != nil {
		return err
	}

	queryParams := &types.QueryRoomParams{
		HotelID: oid,
	}

	rooms, err := h.store.RoomStore.GetRooms(c.Context(), queryParams)
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}
