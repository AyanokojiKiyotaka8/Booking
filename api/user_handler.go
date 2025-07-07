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

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return NewError(http.StatusBadRequest, "invalid user ID format")
	}

	filter := bson.M{"_id": oid}
	user, err := h.userStore.GetUser(c.Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound()
		}
		return ErrInternalServerError()
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	filter := bson.M{}
	users, err := h.userStore.GetUsers(c.Context(), filter)
	if err != nil {
		return ErrInternalServerError()
	}
	return c.JSON(users)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var userParams types.CreateUserParams
	if err := c.BodyParser(&userParams); err != nil {
		return ErrBadRequest()
	}
	if errors := userParams.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}

	user, err := types.NewUserFromParams(&userParams)
	if err != nil {
		return ErrBadRequest()
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return ErrInternalServerError()
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return NewError(http.StatusBadRequest, "invalid user ID format")
	}

	filter := bson.M{"_id": oid}
	err = h.userStore.DeleteUser(c.Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound()
		}
		return ErrInternalServerError()
	}
	return c.JSON(map[string]string{"deleted": id})
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return NewError(http.StatusBadRequest, "invalid user ID format")
	}

	var params types.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": params.ToBSON(),
	}
	err = h.userStore.UpdateUser(c.Context(), filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound()
		}
		return ErrInternalServerError()
	}
	return c.JSON(map[string]string{"updated": id})
}
