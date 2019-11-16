package models

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Application struct {
	ApplicationID primitive.ObjectID `bson:"_id"`
	Status        string             `bson:"status"`
	Job           primitive.ObjectID `bson:"job"`
	Applicant     primitive.ObjectID `bson:"applicant"`
}

func (db *DB) GetApplicationByApplicantAndJob(application *Application, job primitive.ObjectID, applicant primitive.ObjectID) error {
	collection := db.Database(viper.GetString("mongo_db")).Collection("applications")
	err := collection.FindOne(context.Background(),
		bson.M{"job": job, "applicant": applicant}).Decode(application)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetApplicationsByJob(applications *[]Application, jobID string) error {
	job, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return err
	}
	collection := db.Database(viper.GetString("mongo_db")).Collection("applications")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.M{"job": job})
	if err != nil {
		return err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var temp Application
		err := cur.Decode(&temp)
		if err != nil {
			return err
		}
		*applications = append(*applications, temp)
	}
	return nil
}

func (db *DB) CreateApplication(applicant primitive.ObjectID, job primitive.ObjectID, status int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("applications")
	res, err := collection.InsertOne(
		ctx, bson.M{"applicant": applicant, "job": job, "status": status})
	if err != nil {
		return "", errors.Wrap(err, "insert application failed")
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("could not convert to string")
	}
	return id.Hex(), nil
}

type UpdateApplicationPayload struct {
	ApplicationID primitive.ObjectID
	Status        string
}

func (db *DB) UpdateApplication(application UpdateApplicationPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("applications")
	res, err := collection.UpdateOne(
		ctx, bson.M{"_id": application.ApplicationID}, bson.M{"status": application.Status})
	if err != nil {
		return errors.Wrap(err, "update application failed")
	}
	count := res.ModifiedCount
	if count != 1 {
		return errors.Wrap(err, "modify count wrong")
	}
	return nil
}
