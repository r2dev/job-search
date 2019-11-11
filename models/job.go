package models

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type location struct {
	Type        string
	Coordinates [2]int
}
type address struct {
	UnitNumber string
	City       string
	State      string
	Country    string
	Location   location
}

type Job struct {
	JobID         primitive.ObjectID `bson:"_id"`
	Title         string             `bson:"title"`
	Category      string             `bson:"category"`
	Type          string             `bson:"type"`
	Address       address
	FirstSalary   float64            `bson:"firstSalary"`
	SecondSalary  float64            `bson:"secondSalary"`
	PaymentMethod string             `bson:"paymentMethod"`
	Currency      string             `bson:"currency"`
	Rate          string             `bson:"rate"`
	StartDate     time.Time          `bson:"startDate"`
	EndDate       time.Time          `bson:"endDate"`
	StartTime     time.Time          `bson:"startTime"`
	EndTime       time.Time          `bson:"endTime"`
	Description   string             `bson:"description"`
	Reminder      string             `bson:"reminder"`
	Company       primitive.ObjectID `bson:"company"`
	Creator       primitive.ObjectID `bson:"creator"`
}

type CreateJobPayload struct {
	Title         string
	Type          string
	Category      string
	FirstSalary   float64
	SecondSalary  float64
	Currency      string
	Rate          string
	PaymentMethod string
	StartDate     time.Time
	EndDate       time.Time
	StartTime     time.Time
	EndTime       time.Time
	Description   string
	Company       primitive.ObjectID
	Reminder      string
	Creator       primitive.ObjectID
}

func (db *DB) CreateJob(job *CreateJobPayload) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("jobs")
	res, err := collection.InsertOne(ctx,
		bson.M{"title": job.Title, "type": job.Type, "category": job.Category,
			"firstSalary": job.FirstSalary, "secondSalary": job.SecondSalary, "currency": job.Currency, "rate": job.Rate, "paymentMethod": job.PaymentMethod,
			"startDate": job.StartDate, "endDate": job.EndDate, "startTime": job.StartTime, "endTime": job.EndTime,
			"description": job.Description, "company": job.Company, "creator": job.Creator,
		})
	if err != nil {
		return "", errors.Wrap(err, "insert job failed")
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("could not convert to string")
	}
	return id.Hex(), nil
}

func (db *DB) GetJobByID(id string) (*Job, error) {
	var result *Job
	collection := db.Database(viper.GetString("mongo_db")).Collection("jobs")
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
