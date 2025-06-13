package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AyanokojiKiyotaka8/Booking/db"
	"github.com/AyanokojiKiyotaka8/Booking/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME)

	hotel := types.Hotel{
		Name:     "Aaa",
		Location: "Bbb",
	}

	rooms := []types.Room{
		{
			NumPersons: 1,
			Price:      1.11,
		},
		{
			NumPersons: 2,
			Price:      2.22,
		},
		{
			NumPersons: 3,
			Price:      3.33,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertedHotel)

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(context.Background(), &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}
}
