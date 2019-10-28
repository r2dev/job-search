package helpers

import (
	"context"
	"errors"

	"github.com/go-chi/jwtauth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrorNotExisted = errors.New("no user id found at jwt")
var ErrorConvertString = errors.New("id not able to convert to string")
var ErrorConvertObjectID = errors.New("convert to objectID failed")

func GetUserIDFromJWT(context context.Context) (string, primitive.ObjectID, error) {
	_, claims, _ := jwtauth.FromContext(context)
	claimsUser, ok := claims["user_id"]
	if !ok {
		return "", primitive.ObjectID{}, ErrorNotExisted
	}
	userID, ok := claimsUser.(string)
	if !ok {
		return "", primitive.ObjectID{}, ErrorConvertString
	}
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", primitive.ObjectID{}, ErrorConvertObjectID
	}
	return userID, userObjectID, nil
}
