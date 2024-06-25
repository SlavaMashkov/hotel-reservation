package api

import (
	"github.com/SlavaMashkov/hotel-reservation/db"
	"github.com/SlavaMashkov/hotel-reservation/types"
	"github.com/SlavaMashkov/hotel-reservation/utility"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hotelStore db.HotelStore, roomStore db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hotelStore,
		roomStore:  roomStore,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var queryParams types.QueryHotelParams
	if err := c.QueryParser(&queryParams); err != nil {
		return err
	}

	slog.Info("query params", slog.Any("value", queryParams))

	hotels, err := h.hotelStore.GetHotels(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(hotels)
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

	rooms, err := h.roomStore.GetRooms(c.Context(), queryParams)
	if err != nil {
		return err
	}

	if len(rooms) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "rooms not found",
			"id":      id,
		})
	}

	return c.JSON(rooms)
}
