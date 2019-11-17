package models

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	EventInterview = iota + 1
	EventOffer
	EventWork
)

type Event struct {
	EventType   int
	Application primitive.ObjectID
	HireManager primitive.ObjectID
	TimeOptions []time.Time
}

func (db *DB) CreateInterviewEvent(
	application primitive.ObjectID,
	hireManager primitive.ObjectID,
	timeOptions ...time.Time,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("events")
	var timeOptionsBson bson.A
	for _, timeOption := range timeOptions {
		timeOptionsBson = append(timeOptionsBson, primitive.NewDateTimeFromTime(timeOption))
	}
	res, err := collection.InsertOne(ctx,
		bson.M{"application": application, "hireManager": hireManager, "timeOptions": timeOptionsBson})
	if err != nil {
		return errors.Wrap(err, "insert event failed")
	}
	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("could not convert to string")
	}
	return nil
}
