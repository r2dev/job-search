package models

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	UserID   primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Phone    string             `bson:"phone"`
	Email    string             `bson:"email"`
}

var NoFoundUser = errors.New("user no found")

func (db *DB) GetUserByUsername(username string) (User, error) {
	var result User
	collection := db.Database(viper.GetString("mongo_db")).Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, NoFoundUser
		}
		return User{}, err
	}
	return result, nil
}

func (db *DB) GetUserByPhoneNumber(phone string) (User, error) {
	var result User
	collection := db.Database(viper.GetString("mongo_db")).Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"phone": phone}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, nil
		}
		return User{}, err
	}
	return result, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	var result User
	collection := db.Database(viper.GetString("mongo_db")).Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, nil
		}
		return User{}, err
	}
	return result, nil
}

func (db *DB) CreateUserWithUsernameAndPassword(username string, password string) (string, error) {
	collection := db.Database("demo").Collection("users")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate hash password")
	}
	res, err := collection.InsertOne(context.Background(), bson.M{"username": username, "password": string(hashedPassword)})
	if err != nil {
		return "", errors.Wrap(err, "failed to insert user")
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("could not convert to string")
	}
	return id.Hex(), nil
}

func (db *DB) CreateUserWithPhone(phone string) (string, error) {
	collection := db.Database("demo").Collection("users")
	res, err := collection.InsertOne(context.Background(), bson.M{"phone": phone, "verified": false})
	if err != nil {
		return "", errors.Wrap(err, "failed to insert user")
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("could not convert to string")
	}
	return id.Hex(), nil
}
