package api

import (
	"context"
	"log"
	"testing"

	"github.com/AyanokojiKiyotaka8/Booking/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	store  *db.Store
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	store := &db.Store{
		User:    db.NewMongoUserStore(client, db.TestDBNAME),
		Hotel:   db.NewMongoHotelStore(client, db.TestDBNAME),
		Booking: db.NewMongoBookingStore(client, db.TestDBNAME),
	}
	store.Room = db.NewMongoRoomStore(client, db.TestDBNAME, store.Hotel)
	return &testdb{
		client: client,
		store:  store,
	}
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.TestDBNAME).Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
}
