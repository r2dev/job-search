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
type GetApplicationPayload struct {
	Job       primitive.ObjectID
	Applicant primitive.ObjectID
}

func GetApplication(application *GetApplicationPayload) (Application, error) {
	var result Application
	collection := client.Database(viper.GetString("mongo_db")).Collection("applications")
	err := collection.FindOne(context.Background(),
		bson.M{"job": application.Job, "applicant": application.Applicant}).Decode(&result)
	if err != nil {
		return Application{}, err
	}
	return result, nil
}

type CreateApplicationPayload struct {
	Applicant primitive.ObjectID
	Job       primitive.ObjectID
}

func CreateApplication(application *CreateApplicationPayload) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(viper.GetString("mongo_db")).Collection("applications")
	res, err := collection.InsertOne(
		ctx, bson.M{"applicant": application.Applicant, "job": application.Job, "status": "applying"})
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

func UpdateApplication(application UpdateApplicationPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(viper.GetString("mongo_db")).Collection("applications")
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
