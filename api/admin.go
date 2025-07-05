package api

import (
	"github.com/AyanokojiKiyotaka8/Booking/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok || !user.IsAdmin {
		return ErrForbidden()
	}
	return c.Next()
}
