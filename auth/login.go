package auth

import (
	"encoding/json"
	"hirine/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tj/go/http/response"
	"golang.org/x/crypto/bcrypt"
)

type LoginWithPasswordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginWithPassword(w http.ResponseWriter, r *http.Request) {
	var request LoginWithPasswordRequest
	decorder := json.NewDecoder(r.Body)
	err := decorder.Decode(&request)
	if err != nil {
		log.WithError(err).Info("decode failed")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	username := request.Username
	password := request.Password
	user, err := models.GetUserByUsername(username)
	if err != nil {
		if err == models.NoFoundUser {
			log.WithError(err).Info("dont get user")
		}
		response.Error(w, http.StatusBadRequest)
		return
	}
	encryptedPassword := user.Password
	err = bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.WithError(err).Info("mismatch hash")
			response.Error(w, http.StatusBadRequest)
			return
		} else {
			log.WithError(err).Info("compare hash unknwon error")
			response.Error(w, http.StatusInternalServerError)
			return
		}
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.UserID.Hex()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(viper.GetString("jwt_secret")))
	if err != nil {
		log.WithError(err).Warn("sign error")
		response.Error(w, http.StatusInternalServerError)
		return
	}
	response.OK(w, map[string]string{
		"token": t,
	})

}
