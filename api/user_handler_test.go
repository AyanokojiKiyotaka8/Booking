package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/AyanokojiKiyotaka8/Booking/types"
	"github.com/gofiber/fiber/v2"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})
	userHandler := NewUserHandler(tdb.store.User)
	app.Post("/", userHandler.HandlePostUser)

	params := &types.CreateUserParams{
		FirstName: "qqq",
		LastName:  "www",
		Email:     "qqq@www.com",
		Password:  "qqww123",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var user *types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user is nil after decoding response")
	}

	if user.ID.IsZero() {
		t.Errorf("expected user ID to be set")
	}
	if len(user.EncPassword) > 0 {
		t.Errorf("Encrypted password should not be included in JSON")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected first name %s, got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected last name %s, got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s, got %s", params.Email, user.Email)
	}
}
