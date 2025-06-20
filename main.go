package main

import (
	"context"
	"flag"
	"log"

	"github.com/AyanokojiKiyotaka8/Booking/api"
	"github.com/AyanokojiKiyotaka8/Booking/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// storages
	userStore := db.NewMongoUserStore(client, db.DBNAME)
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME, hotelStore)

	// handlers
	userHandler := api.NewUserHandler(userStore)
	hotelHandler := api.NewHotelHandler(hotelStore, roomStore)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.JSON(map[string]string{"error": err.Error()})
		},
	})
	apiv1 := app.Group("/api/v1")

	// user API's
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	// hotel API's
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)

	app.Listen(*listenAddr)
}
