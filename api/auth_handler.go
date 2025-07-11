package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AyanokojiKiyotaka8/Booking/db"
	"github.com/AyanokojiKiyotaka8/Booking/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	filter := bson.M{"email": params.Email}
	user, err := h.userStore.GetUser(c.Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUnauthorized()
		}
		return ErrInternalServerError()
	}

	if !types.IsValidPassword(user.EncPassword, params.Password) {
		return ErrUnauthorized()
	}

	resp := &AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}
	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {
	claims := jwt.MapClaims{
		"id":    user.ID.Hex(),
		"email": user.Email,
		"exp":   jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("problems with token", err)
	}
	return tokenStr
}
