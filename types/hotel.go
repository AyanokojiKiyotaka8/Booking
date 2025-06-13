package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
}

type Room struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NumPersons int64              `bson:"numPersons" json:"numPersons"`
	Price      float64            `bson:"price" json:"price"`
	HotelID    primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}
