package main

import (
	"context"
	"log"
	"os"

	"github.com/AyanokojiKiyotaka8/Booking/api"
	"github.com/AyanokojiKiyotaka8/Booking/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// DBs
	store := &db.Store{
		User:    db.NewMongoUserStore(client, db.DBNAME),
		Hotel:   db.NewMongoHotelStore(client, db.DBNAME),
		Booking: db.NewMongoBookingStore(client, db.DBNAME),
	}
	store.Room = db.NewMongoRoomStore(client, db.DBNAME, store.Hotel)

	// handlers
	userHandler := api.NewUserHandler(store.User)
	hotelHandler := api.NewHotelHandler(store)
	authHandler := api.NewAuthHandler(store.User)
	roomHandler := api.NewRoomHandler(store)
	bookingHandler := api.NewBookingHandler(store)

	app := fiber.New(fiber.Config{
		ErrorHandler: api.ErrorHandler,
	})
	apiv1 := app.Group("/api/v1", api.JWTAuthentication(store.User))
	admin := apiv1.Group("/admin", api.AdminAuth)
	auth := app.Group("/api")

	// auth APIs
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// user APIs
	admin.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	// hotel APIs
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// room APIs
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiv1.Get("/room", roomHandler.HandleGetRooms)

	// booking APIs with admin authorized
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	// booking APIs with user authorized
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	app.Listen(os.Getenv("HTTP_LISTEN_ADDRESS"))
}

func init() {
	db.LoadConfig()
}
