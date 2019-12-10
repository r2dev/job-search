package app

import (
	"encoding/json"
	"github.com/r2dev/job-search/helpers"
	"github.com/r2dev/job-search/models"
	"net/http"

	"github.com/tj/go/http/response"
)

type CreateCompanyRequest struct {
	CompanyName string `json:"companyName"`
}

func (app *App) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var request CreateCompanyRequest
	decorder := json.NewDecoder(r.Body)
	err := decorder.Decode(&request)
	if err != nil {
		app.L.WithError(err).Info("decode failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	_, userObjectID, err := helpers.GetUserIDFromJWT(r.Context())
	if err != nil {
		app.L.WithError(err).Info(" get id failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}

	companyName := request.CompanyName
	id, err := app.DB.CreateCompany(&models.CreateCompanyPayload{
		CompanyName: companyName,
		Admin:       userObjectID,
	})
	if err != nil {
		app.L.WithError(err).Info("create company failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}

	response.OK(w, map[string]string{
		"id": id,
	})
}

func (app *App) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	// @todo delete
}

func (app *App) GetCompany(w http.ResponseWriter, r *http.Request) {
	response.OK(w, 1)
}

func (app *App) UpdateCompany(w http.ResponseWriter, r *http.Request) {

}
