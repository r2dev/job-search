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
	EventID     primitive.ObjectID
	EventType   int
	EventTime   time.Time
	Application primitive.ObjectID
	HireManager primitive.ObjectID
	Applicant   primitive.ObjectID
	TimeOptions []time.Time
	Confirmed   bool
	Status      int
}

type StatusInterview int

const (
	StatusInterviewCreated StatusInterview = iota + 1
	StatusInterviewUpdated
	StatusInterviewConfirmed
	StatusInterviewDeclined
	StatusInterviewCancel
)

func (db *DB) CreateInterviewEvent(
	application primitive.ObjectID,
	candidate primitive.ObjectID,
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
		bson.M{"application": application, "hireManager": hireManager,
			"timeOptions": timeOptionsBson, "status": StatusInterviewCreated, "candidate": candidate})
	if err != nil {
		return errors.Wrap(err, "insert event failed")
	}
	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("could not convert to string")
	}
	return nil
}

func (db *DB) ConfirmInterviewEvent(eventID primitive.ObjectID, timeOption time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("events")
	res, err := collection.UpdateOne(ctx, bson.M{"_id": eventID},
		bson.M{"status": StatusInterviewConfirmed, "eventTime": primitive.NewDateTimeFromTime(timeOption)})
	if err != nil {
		return errors.Wrap(err, "confirm interview event failed")
	}
	count := res.ModifiedCount
	if count != 1 {
		return errors.Wrap(err, "modify count wrong")
	}
	return nil
}

func (db *DB) DeclineInterviewEvent(eventID primitive.ObjectID, reason string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("events")
	res, err := collection.UpdateOne(ctx, bson.M{"_id": eventID},
		bson.M{"reason": reason, "status": StatusInterviewDeclined})
	if err != nil {
		return errors.Wrap(err, "decline interview event failed")
	}
	count := res.ModifiedCount
	if count != 1 {
		return errors.Wrap(err, "modify count wrong")
	}
	return nil
}

func (db *DB) UpdateInterviewEvent(eventID primitive.ObjectID, hireManager primitive.ObjectID, timeOptions ...time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("events")
	var timeOptionsBson bson.A
	for _, timeOption := range timeOptions {
		timeOptionsBson = append(timeOptionsBson, primitive.NewDateTimeFromTime(timeOption))
	}
	res, err := collection.UpdateOne(ctx, bson.M{"_id": eventID},
		bson.M{"hireManger": hireManager, "timeOptions": timeOptions, "status": StatusInterviewUpdated})
	if err != nil {
		return errors.Wrap(err, "decline interview event failed")
	}
	count := res.ModifiedCount
	if count != 1 {
		return errors.Wrap(err, "modify count wrong")
	}
	return nil
}

func (db *DB) CancelInterviewEvent(eventID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("events")

	res, err := collection.UpdateOne(ctx, bson.M{"_id": eventID},
		bson.M{"status": StatusInterviewCancel})
	if err != nil {
		return errors.Wrap(err, "decline interview event failed")
	}
	count := res.ModifiedCount
	if count != 1 {
		return errors.Wrap(err, "modify count wrong")
	}
	return nil
}

func (db *DB) GetEventByEventID(event *Event, eventIDString string) error {
	eventID, err := primitive.ObjectIDFromHex(eventIDString)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("events")
	err = collection.FindOne(ctx, bson.M{"_id": eventID}).Decode(event)
	if err != nil {
		return err
	}
	return nil
}
