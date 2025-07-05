package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AyanokojiKiyotaka8/Booking/db/fixtures"
	"github.com/AyanokojiKiyotaka8/Booking/types"
	"github.com/gofiber/fiber/v2"
)

func TestGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.store, false, "qwe", "rty")
	otherUser := fixtures.AddUser(tdb.store, false, "aaa", "bbb")
	hotel := fixtures.AddHotel(tdb.store, "qqq", "www", 3)
	room := fixtures.AddRoom(tdb.store, "large", true, 123.45, hotel.ID)
	booking := fixtures.AddBoking(tdb.store, room.ID, user.ID, 7, time.Now(), time.Now().AddDate(0, 0, 3))

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})
	auth := app.Group("/", JWTAuthentication(tdb.store.User))
	bookingHandler := NewBookingHandler(tdb.store)
	auth.Get("/:id", bookingHandler.HandleGetBooking)

	t.Run("user can fetch their own booking", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var bookingResp *types.Booking
		if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
			t.Fatal(err)
		}

		if booking.ID != bookingResp.ID {
			t.Errorf("expected booking ID %v, got %v", booking.ID, bookingResp.ID)
		}
		if booking.UserID != bookingResp.UserID {
			t.Errorf("expected user ID %v, got %v", booking.UserID, bookingResp.UserID)
		}
	})

	t.Run("other user cannot access booking", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Api-Token", CreateTokenFromUser(otherUser))

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("expected status 403, got %d", resp.StatusCode)
		}
	})
}

func TestGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.store, false, "qwe", "rty")
	admin := fixtures.AddUser(tdb.store, true, "admin", "admin")
	hotel := fixtures.AddHotel(tdb.store, "qqq", "www", 3)
	room := fixtures.AddRoom(tdb.store, "large", true, 123.45, hotel.ID)
	booking := fixtures.AddBoking(tdb.store, room.ID, user.ID, 7, time.Now(), time.Now().AddDate(0, 0, 3))

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})
	auth := app.Group("/", JWTAuthentication(tdb.store.User), AdminAuth)
	bookingHandler := NewBookingHandler(tdb.store)
	auth.Get("/", bookingHandler.HandleGetBookings)

	t.Run("admin can fetch all bookings", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Api-Token", CreateTokenFromUser(admin))

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var bookings []*types.Booking
		if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
			t.Fatal(err)
		}
		if len(bookings) != 1 {
			t.Fatalf("expected 1 booking, got %d", len(bookings))
		}
		if booking.ID != bookings[0].ID {
			t.Errorf("expected booking ID %v, got %v", booking.ID, bookings[0].ID)
		}
	})

	t.Run("non-admin cannot fetch all bookings", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusForbidden {
			t.Fatalf("expected 401 or 403, got %d", resp.StatusCode)
		}
	})
}
