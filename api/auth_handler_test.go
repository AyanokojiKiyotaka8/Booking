package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/AyanokojiKiyotaka8/Booking/db"
	"github.com/AyanokojiKiyotaka8/Booking/types"
	"github.com/gofiber/fiber/v2"
)

func makeTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(&types.CreateUserParams{
		FirstName: "qwe",
		LastName:  "rty",
		Email:     "qwe@rty.com",
		Password:  "qwerty123",
	})

	if err != nil {
		t.Fatal(err)
	}

	insertedUser, err := userStore.InsertUser(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}
	return insertedUser
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := makeTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "qwe@rty.com",
		Password: "qwerty123",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected to get http status code 200, but got %d", resp.StatusCode)
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
		t.Fatal("expected the user to be the inserted user")
	}
}

func TestAuthenticateFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	makeTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "qwe@rty.com",
		Password: "qwerty12345",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected to get status code 400, but got %d", resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("expected the type of gen response to be error, but got %s", genResp.Type)
	}

	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected the message of gen response to be <invalid credentials>, but got %s", genResp.Msg)
	}
}
