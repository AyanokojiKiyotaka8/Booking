package main

import (
	"context"
	"log"

	"github.com/AyanokojiKiyotaka8/Booking/db"
	"github.com/AyanokojiKiyotaka8/Booking/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	userStore  db.UserStore
)

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:    "small",
			Price:   1.11,
			SeaSide: false,
		},
		{
			Size:    "normal",
			Price:   2.22,
			SeaSide: false,
		},
		{
			Size:    "large",
			Price:   3.33,
			SeaSide: true,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(context.Background(), &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func seedUser(isAdmin bool, fname, lname, email, password string) {
	user, err := types.NewUserFromParams(&types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})

	user.IsAdmin = isAdmin

	if err != nil {
		log.Fatal(err)
	}

	_, err = userStore.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	seedHotel("Aaa", "Bbb", 1)
	seedHotel("Qqq", "Www", 2)
	seedHotel("Eee", "Rrr", 3)
	seedUser(false, "qwe", "rty", "qwe@rty.com", "qwerty123")
	seedUser(true, "admin", "admin", "admin@admin.com", "admin123")
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
	roomStore = db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
	userStore = db.NewMongoUserStore(client, db.DBNAME)
}
