package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	EventInterview = iota + 1
	EventWork
)

type Event struct {
	EventID        primitive.ObjectID `bson:"_id"`
	EventType      int                `bson:"eventType"`
	EventTime      time.Time          `bson:"eventTime"`
	ActionRequired bool               `bson:"actionRequired"`
	Status         int                `bson:"status"`
	Attendee       primitive.ObjectID `bson:"Attendee"`

	// interview event specific field
	Application primitive.ObjectID `bson:"application"`
	HireManager primitive.ObjectID `bson:"hireManager"`
	TimeOptions []time.Time        `bson:"timeOptions"`
}

type StatusInterview int

const (
	StatusInterviewCreated StatusInterview = iota + 1
	StatusInterviewUpdated
	StatusInterviewConfirmed
	StatusInterviewDeclined
	StatusInterviewCancel

	StatusWorkCreated
	StatusWorkConfirmed
	StatusWorkDeclined
	StatusWorkCancel
)

func (db *DB) CreateInterviewEvent(
	application primitive.ObjectID,
	attendee primitive.ObjectID,
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
		bson.M{"application": application, "hireManager": hireManager, "timeOptions": timeOptionsBson, "status": StatusInterviewCreated, "attendee": attendee, "eventType": EventInterview})
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
		bson.M{"$set": bson.M{"status": StatusInterviewConfirmed, "eventTime": primitive.NewDateTimeFromTime(timeOption)}})
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

func (db *DB) UpdateInterviewEvent(eventID primitive.ObjectID, hireManager primitive.ObjectID, timeOptions []time.Time) error {
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

	res, err := collection.UpdateOne(ctx, bson.M{"_id": eventID}, bson.M{"status": StatusInterviewCancel})
	if err != nil {
		return errors.Wrap(err, "decline interview event failed")
	}
	count := res.ModifiedCount
	if count != 1 {
		return errors.Wrap(err, "modify count wrong")
	}
	return nil
}

func (db *DB) CreateWorkEvent(attendee primitive.ObjectID, eventTime time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("events")
	res, err := collection.InsertOne(ctx,
		bson.M{"eventTime": eventTime, "status": StatusWorkCreated, "attendee": attendee, "eventType": EventWork})
	if err != nil {
		return errors.Wrap(err, "insert event failed")
	}
	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("could not convert to string")
	}
	return nil
}

func (db *DB) ConfirmWorkEvent(eventID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("events")
	res, err := collection.UpdateOne(ctx, bson.M{"_id": eventID},
		bson.M{"$set": bson.M{"status": StatusWorkConfirmed}})
	if err != nil {
		return errors.Wrap(err, "confirm work event failed")
	}
	count := res.ModifiedCount
	if count != 1 {
		return errors.Wrap(err, "modify count wrong")
	}
	return nil
}

func (db *DB) DeclineWorkEvent(eventID primitive.ObjectID, reason string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("events")
	res, err := collection.UpdateOne(ctx, bson.M{"_id": eventID},
		bson.M{"$set": bson.M{"reason": reason, "status": StatusWorkDeclined}})
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

func (db *DB) GetEventsByAttendee(events *[]Event, attendee primitive.ObjectID, limit int, skip int) error {
	collection := db.Database(viper.GetString("mongo_db")).Collection("events")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.M{"attendee": attendee}, options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)))
	if err != nil {
		return err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var temp Event
		err := cur.Decode(&temp)
		if err != nil {
			return err
		}
		*events = append(*events, temp)
	}
	return nil
}
