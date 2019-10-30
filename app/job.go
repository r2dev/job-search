package app

import (
	"hirine/helpers"
	"hirine/models"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tj/go/http/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateJobRequest struct {
	Title           string `json:"title"`
	Category        string `json:"category"`
	Type            string `json:"type"`
	FirstSalary     int64  `json:"firstSalary"`
	SecondSalary    int64  `json:"secondSalary"`
	PaymentMethod   string `json:"paymentMethod"`
	Currency        string `json:"currency"`
	Rate            string `json:"rate"`
	StartDateString string `json:"startDate"`
	EndDateString   string `json:"endDate"`
	StartTimeString string `json:"startTime"`
	EndTimeString   string `json:"endTime"`
	Description     string `json:"description"`
	Reminder        string `json:"reminder"`
	Company         string `json:"company"`
}

func (app *App) CreateJob(w http.ResponseWriter, r *http.Request) {
	var request CreateJobRequest
	err := helpers.DecodeJSON(r, &request)
	if err != nil {
		log.WithError(err).Info("decode failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	_, userObjectID, err := helpers.GetUserIDFromJWT(r.Context())
	if err != nil {
		log.WithError(err).Info("get id failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	title := request.Title
	category := request.Category
	firstSalary := request.FirstSalary
	secondSalary := request.SecondSalary
	paymentMethod := request.PaymentMethod
	currency := request.Currency
	rate := request.Rate
	startDate := helpers.ParseJavascriptTimeString(request.StartDateString)
	endDate := helpers.ParseJavascriptTimeString(request.EndDateString)
	startTime := helpers.ParseJavascriptTimeString(request.StartTimeString)
	endTime := helpers.ParseJavascriptTimeString(request.EndTimeString)
	description := request.Description
	reminder := request.Reminder
	company := request.Company

	companyObjectID, err := primitive.ObjectIDFromHex(company)
	if err != nil {
		log.WithError(err).Info("convert to objectID failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	companyValue, err := app.DB.GetCompanyById(companyObjectID)
	if err != nil {
		log.WithError(err).Info("get company failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	if companyValue.Admin != userObjectID {
		log.WithError(err).Info("unauthorized")
		response.Error(w, http.StatusUnauthorized)
		return
	}
	id, err := app.DB.CreateJob(&models.CreateJobPayload{
		Title:         title,
		Category:      category,
		FirstSalary:   firstSalary,
		SecondSalary:  secondSalary,
		PaymentMethod: paymentMethod,
		Currency:      currency,
		Rate:          rate,
		StartDate:     startDate,
		EndDate:       endDate,
		StartTime:     startTime,
		EndTime:       endTime,
		Description:   description,
		Reminder:      reminder,
		Company:       companyObjectID,
		Creator:       userObjectID,
	})
	if err != nil {
		log.WithError(err).Info("create job failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}

	response.OK(w, map[string]string{
		"id": id,
	})

	// create cache
}

func (app *App) DeleteJob(w http.ResponseWriter, r *http.Request) {
	// update cache
}

func (app *App) UpdateJob(w http.ResponseWriter, r *http.Request) {
	// update cache
}
