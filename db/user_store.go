package db

import (
	"context"

	"github.com/AyanokojiKiyotaka8/Booking/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type UserStore interface {
	GetUser(context.Context, bson.M) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userColl),
	}
}

func (s *MongoUserStore) GetUser(ctx context.Context, filter bson.M) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
