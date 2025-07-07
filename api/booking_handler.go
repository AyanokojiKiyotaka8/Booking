package api

import (
	"errors"
	"net/http"

	"github.com/AyanokojiKiyotaka8/Booking/db"
	"github.com/AyanokojiKiyotaka8/Booking/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	filter := bson.M{}
	bookings, err := h.store.Booking.GetBookings(c.Context(), filter)
	if err != nil {
		return ErrInternalServerError()
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return NewError(http.StatusBadRequest, "invalid booking ID format")
	}

	filter := bson.M{"_id": id}
	booking, err := h.store.Booking.GetBooking(c.Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound()
		}
		return ErrInternalServerError()
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrInternalServerError()
	}

	if booking.UserID != user.ID {
		return ErrForbidden()
	}
	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return NewError(http.StatusBadRequest, "invalid booking ID format")
	}

	filter := bson.M{"_id": id}
	booking, err := h.store.Booking.GetBooking(c.Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound()
		}
		return ErrInternalServerError()
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrInternalServerError()
	}

	if !user.IsAdmin && booking.UserID != user.ID {
		return ErrForbidden()
	}

	filter = bson.M{"_id": booking.ID}
	update := bson.M{
		"$set": bson.M{
			"cancelled": true,
		},
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), filter, update); err != nil {
		return ErrInternalServerError()
	}
	return c.JSON(fiber.Map{"cancelled": booking.ID.Hex()})
}
