package api

import (
	"fmt"
	"os"

	"github.com/AyanokojiKiyotaka8/Booking/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-Api-Token")
		if len(token) == 0 {
			return ErrUnauthorized()
		}

		claims, err := validateToken(token)
		if err != nil {
			return ErrUnauthorized()
		}

		idStr, ok := claims["id"].(string)
		if !ok {
			return ErrUnauthorized()
		}

		userID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return ErrUnauthorized()
		}

		filter := bson.M{"_id": userID}
		user, err := userStore.GetUser(c.Context(), filter)
		if err != nil {
			return ErrUnauthorized()
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, ErrUnauthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse JWT token", err)
		return nil, ErrUnauthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		fmt.Println("token is invalid")
		return nil, ErrUnauthorized()
	}
	return claims, nil
}
