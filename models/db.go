package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	*mongo.Client
}

func InitMongo(mongoUrl string) (*DB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		// log.Panic(err)
		return nil, err
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		// log.Panic(err)
		return nil, err
	}
	return &DB{client}, nil
}
