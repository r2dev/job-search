package models

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
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
	Title         string
	Address       address
	Category      string
	Type          string
	PaymentMethod string
	startDate     time.Time
	endDate       time.Time
	startTime     time.Time
	endTime       time.Time
	description   string
	reminder      string
}

func Create(job *Job) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("demo").Collection("jobs")
	res, err := collection.InsertOne(ctx, bson.M{"title": job.Title, "category": job.Category})
	if err != nil {
		return "", errors.Wrap(err, "insert job failed")
	}
	id, ok := res.InsertedID.(string)
	if !ok {
		return "", errors.New("could not convert to string")
	}
	return id, nil
}
