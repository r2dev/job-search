package company

import (
	"encoding/json"
	"hirine/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-chi/jwtauth"

	log "github.com/sirupsen/logrus"
	"github.com/tj/go/http/response"
)

type CreateCompanyRequest struct {
	CompanyName string `json:"companyName"`
}

func CreateCompany(w http.ResponseWriter, r *http.Request) {
	var request CreateCompanyRequest
	decorder := json.NewDecoder(r.Body)
	err := decorder.Decode(&request)
	if err != nil {
		log.WithError(err).Info("decode failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	claimsUser, ok := claims["user_id"]
	if !ok {
		log.WithError(err).Info("jwt user_id not found")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	userID, ok := claimsUser.(string)
	if !ok {
		log.WithError(err).Info("claims convert failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.WithError(err).Info(" convert to objectID failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}

	companyName := request.CompanyName
	id, err := models.CreateCompany(models.CreateCompanyPayload{
		CompanyName: companyName,
		Admin:       userObjectID,
	})
	if err != nil {
		log.WithError(err).Info("create company failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}

	response.OK(w, map[string]string{
		"id": id,
	})
}
