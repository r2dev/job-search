package models

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	CompanyID    primitive.ObjectID   `bson:"_id"`
	CompanyName  string               `bson:"companyName"`
	ProfileImage string               `bson:"profileImage"`
	Verify       bool                 `bson:"verify"`
	Admin        primitive.ObjectID   `bson:"admin"`
	Manager      []primitive.ObjectID `bson:"Manager"`
}

type CreateCompanyPayload struct {
	Admin        primitive.ObjectID
	CompanyName  string
	ProfileImage string
}

func (db *DB) CreateCompany(company CreateCompanyPayload) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("companys")
	res, err := collection.InsertOne(
		ctx, bson.M{"companyName": company.CompanyName, "verify": false, "admin": company.Admin})
	if err != nil {
		return "", errors.Wrap(err, "insert company failed")
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("could not convert to string")
	}
	return id.Hex(), nil
}

type UpdateCompanyPayload struct {
	CompanyID    string
	CompanyName  string
	ProfileImage string
}

func (db *DB) UpdateCompany(company UpdateCompanyPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Database(viper.GetString("mongo_db")).Collection("companys")
	res, err := collection.UpdateOne(
		ctx, bson.M{"_id": company.CompanyID}, bson.M{"companyName": company.CompanyName, "verify": false})
	if err != nil {
		return errors.Wrap(err, "update company failed")
	}
	count := res.ModifiedCount
	if count != 1 {
		return errors.Wrap(err, "modify count wrong")
	}
	return nil
}

func (db *DB) GetCompanyById(CompanyID primitive.ObjectID) (Company, error) {
	var result Company
	collection := db.Database(viper.GetString("mongo_db")).Collection("companys")
	err := collection.FindOne(context.Background(), bson.M{"_id": CompanyID}).Decode(&result)
	return result, err
}
