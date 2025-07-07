package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AyanokojiKiyotaka8/Booking/api"
	"github.com/AyanokojiKiyotaka8/Booking/db"
	"github.com/AyanokojiKiyotaka8/Booking/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DropDB(client *mongo.Client, dbname string) {
	if err := client.Database(dbname).Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Dropped the DB <%s>\n", dbname)
}

func DropCollection(client *mongo.Client, dbname string, collName string) {
	if err := client.Database(dbname).Collection(collName).Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Dropped the Collection <%s> in DB <%s>\n", collName, dbname)
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	DropDB(client, db.DBNAME)
	DropCollection(client, db.DBNAME, "users")
	DropCollection(client, db.DBNAME, "hotels")
	DropCollection(client, db.DBNAME, "rooms")
	DropCollection(client, db.DBNAME, "bookings")

	store := &db.Store{
		User:    db.NewMongoUserStore(client, db.DBNAME),
		Hotel:   db.NewMongoHotelStore(client, db.DBNAME),
		Booking: db.NewMongoBookingStore(client, db.DBNAME),
	}
	store.Room = db.NewMongoRoomStore(client, db.DBNAME, store.Hotel)

	user := fixtures.AddUser(store, false, "qwe", "rty")
	fmt.Println("user -> ", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, true, "admin", "admin")
	fmt.Println("admin -> ", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "qqq", "www", 3)
	room := fixtures.AddRoom(store, "large", true, 379.99, hotel.ID)
	booking := fixtures.AddBooking(store, room.ID, user.ID, 7, time.Now(), time.Now().AddDate(0, 0, 3))
	fmt.Println(booking)
}

func init() {
	db.LoadConfig()
}
