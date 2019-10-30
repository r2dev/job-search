package app

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

func (app *App) ApplyJob(w http.ResponseWriter, r *http.Request) {
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

	_, err = app.DB.GetApplication(&models.GetApplicationPayload{
		Job:       jobObjectID,
		Applicant: userObjectID,
	})
	if err != nil && err != mongo.ErrNoDocuments {
		log.WithError(err).Info("")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	id, err := app.DB.CreateApplication(&models.CreateApplicationPayload{
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

func (app *App) UpdateApplication(w http.ResponseWriter, r *http.Request) {
	// employer send interview to employee
	// employee accept interview
	// employee decline interview
	// employer cancel interview before interview begin
	// system expire interview
	// system start interview
	// system end interview
	// employee give feedback after the interview about employer
	// employer give feedback after the interview about employee

	// employer send offer to employee
	// employee accept offer
	// employee decline offer
	// employer cancel offer before employee accept offer
	// employer cancel offer after employee accept offer
	// system expire offer

}
