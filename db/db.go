package db

const DBNAME = "booking"
const DBURI = "mongodb://localhost:27017"
const TestDBNAME = "booking-test"

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
