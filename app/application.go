package app

import (
	"encoding/json"
	"hirine/helpers"
	"hirine/models"
	"net/http"

	"github.com/tj/go/http/response"
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
		app.L.WithError(err).Info("decode failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	userID, _, err := helpers.GetUserIDFromJWT(r.Context())
	if err != nil {
		app.L.WithError(err).Info(" get user id failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	var application models.Application
	err = app.DB.GetApplicationByApplicantAndJob(&application, request.Job, userID)
	if err != nil && err != mongo.ErrNoDocuments {
		app.L.WithError(err).Info("")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	id, err := app.DB.CreateApplication(userID, request.Job, StatusApplying)
	if err != nil {
		response.Error(w, http.StatusInternalServerError)
		return
	}

	response.OK(w, map[string]string{
		"id": id,
	})
}
