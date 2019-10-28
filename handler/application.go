package handler

import (
	"encoding/json"
	"hirine/helpers"
	"hirine/models"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tj/go/http/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApplyJobRequest struct {
	Job string `json:"job"`
}

func ApplyJob(w http.ResponseWriter, r *http.Request) {
	var request ApplyJobRequest
	decorder := json.NewDecoder(r.Body)
	err := decorder.Decode(&request)
	if err != nil {
		log.WithError(err).Info("decode failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	_, userObjectID, err := helpers.GetUserIDFromJWT(r.Context())
	if err != nil {
		log.WithError(err).Info(" get user id failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}

	jobObjectID, err := primitive.ObjectIDFromHex(request.Job)
	if err != nil {
		log.WithError(err).Info(" convert to objectID failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}

	_, err = models.GetApplication(&models.GetApplicationPayload{
		Job:       jobObjectID,
		Applicant: userObjectID,
	})
	if err != nil && err != mongo.ErrNoDocuments {
		log.WithError(err).Info("")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	id, err := models.CreateApplication(&models.CreateApplicationPayload{
		Applicant: userObjectID,
		Job:       jobObjectID,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError)
		return
	}

	response.OK(w, map[string]string{
		"id": id,
	})
}
