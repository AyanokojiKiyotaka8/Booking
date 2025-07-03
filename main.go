package main

import (
	"context"
	"flag"
	"log"

	"github.com/AyanokojiKiyotaka8/Booking/api"
	"github.com/AyanokojiKiyotaka8/Booking/api/middleware"
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
	bookingStore := db.NewMongoBookingStore(client, db.DBNAME)
	store := &db.Store{
		User:    userStore,
		Hotel:   hotelStore,
		Room:    roomStore,
		Booking: bookingStore,
	}

	// handlers
	userHandler := api.NewUserHandler(userStore)
	hotelHandler := api.NewHotelHandler(store)
	authHandler := api.NewAuthHandler(userStore)
	roomHandler := api.NewRoomHandler(store)
	bookingHandler := api.NewBookingHandler(store)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.JSON(map[string]string{"error": err.Error()})
		},
	})
	apiv1 := app.Group("/api/v1", middleware.JWTAuthentication(userStore))
	admin := apiv1.Group("/admin", middleware.AdminAuth)
	auth := app.Group("/api")

	// auth API's
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// user API's
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	// hotel API's
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// room API's
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiv1.Get("/room", roomHandler.HandleGetRooms)

	// booking API's with admin authorized
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	// booking API's with user authorized
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	app.Listen(*listenAddr)
}
