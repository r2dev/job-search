package models

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

func InitMongo(mongoUrl string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Panic(err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Panic(err)
	}
}
