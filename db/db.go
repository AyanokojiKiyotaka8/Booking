package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBURI      string
	DBNAME     string
	TestDBNAME string
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

func LoadConfig() {
	// if testing then add "../.env" into Load()
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	DBURI = os.Getenv("MONGO_DB_URI")
	DBNAME = os.Getenv("MONGO_DB_NAME")
	TestDBNAME = os.Getenv("MONGO_TEST_DB_NAME")
}
