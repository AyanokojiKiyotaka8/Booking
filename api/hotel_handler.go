package api

import (
	"errors"
	"net/http"

	"github.com/AyanokojiKiyotaka8/Booking/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	filter := bson.M{}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter)
	if err != nil {
		return ErrInternalServerError()
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return NewError(http.StatusBadRequest, "invalid hotel ID format")
	}

	filter := bson.M{"_id": oid}
	hotel, err := h.store.Hotel.GetHotel(c.Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound()
		}
		return ErrInternalServerError()
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return NewError(http.StatusBadRequest, "invalid hotel ID format")
	}

	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrInternalServerError()
	}
	return c.JSON(rooms)
}
