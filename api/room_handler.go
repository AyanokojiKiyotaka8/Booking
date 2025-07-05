package api

import (
	"net/http"
	"time"

	"github.com/AyanokojiKiyotaka8/Booking/db"
	"github.com/AyanokojiKiyotaka8/Booking/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.Store
}

type RoomBookParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params *RoomBookParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	if time.Now().After(params.FromDate) || !params.TillDate.After(params.FromDate) {
		return NewError(http.StatusBadRequest, "inappropriate booking period")
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrInternalServerError()
	}

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return NewError(http.StatusBadRequest, "invalid room ID format")
	}

	filter := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$lte": params.TillDate,
		},
		"tillDate": bson.M{
			"$gte": params.FromDate,
		},
		"cancelled": false,
	}
	bookings, err := h.store.Booking.GetBookings(c.Context(), filter)
	if err != nil {
		return ErrInternalServerError()
	}
	if len(bookings) > 0 {
		return NewError(http.StatusBadRequest, "room is not available for that period")
	}

	booking := &types.Booking{
		RoomID:     roomID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
		UserID:     user.ID,
	}
	insertedBooking, err := h.store.Booking.InsertBooking(c.Context(), booking)
	if err != nil {
		return ErrInternalServerError()
	}
	return c.JSON(map[string]any{"booked": insertedBooking})
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	filter := bson.M{}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrInternalServerError()
	}
	return c.JSON(rooms)
}
