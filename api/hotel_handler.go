package api

import (
	"github.com/SlavaMashkov/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct {
	hotelStore db.HotelStore
}

func NewHotelHandler(hotelStore db.HotelStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hotelStore,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.hotelStore.GetHotels(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}
