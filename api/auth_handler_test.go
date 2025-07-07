package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/AyanokojiKiyotaka8/Booking/db/fixtures"
	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := fixtures.AddUser(tdb.store, false, "qwe", "rty")

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := &AuthParams{
		Email:    "qwe@rty.com",
		Password: "qwe_rty", // correct password
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected to get HTTP status code 200, but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if authResp.Token == "" {
		t.Fatal("expected the JWT token to be present in auth response")
	}

	insertedUser.EncPassword = ""
	if !reflect.DeepEqual(authResp.User, insertedUser) {
		t.Fatal("expected the user to match inserted user (excluding password)")
	}
}

func TestAuthenticateFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	fixtures.AddUser(tdb.store, false, "qwe", "rty")

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := &AuthParams{
		Email:    "qwe@rty.com",
		Password: "wrong_password", // incorrect password
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status code 401, but got %d", resp.StatusCode)
	}
}
