package handler

import (
	"crypto/rand"
	"encoding/json"
	"hirine/cache"
	"hirine/models"
	"hirine/sms"
	"io"
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

type RegisterWithPasswordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterWithPassword(w http.ResponseWriter, r *http.Request) {
	var request RegisterWithPasswordRequest
	decorder := json.NewDecoder(r.Body)
	err := decorder.Decode(&request)

	username := request.Username
	password := request.Password
	_, err = models.GetUserByUsername(username)
	if err != nil {
		if err != models.NoFoundUser {
			log.WithError(err)
			response.Error(w, http.StatusInternalServerError)
			return
		}
	}
	id, err := models.CreateUserWithUsernameAndPassword(username, password)
	if err != nil {
		log.WithError(err)
		response.Error(w, http.StatusBadGateway)
		return
	}
	response.OK(w, map[string]string{
		"id": id,
	})

}

type RegisterWithPhoneNumberRequest struct {
	Phone string `json:"phone"`
	// Code  string `json:"code"`
}

func RegisterWithPhoneNumber(w http.ResponseWriter, r *http.Request) {
	var request RegisterWithPhoneNumberRequest
	decorder := json.NewDecoder(r.Body)
	err := decorder.Decode(&request)

	phone := request.Phone
	_, err = models.GetUserByPhoneNumber(phone)
	if err != nil {
		log.WithError(err)
		response.Error(w, http.StatusInternalServerError)
		return
	}
	verifyCode := encodeToString(4)
	go cache.Save(phone, verifyCode)
	go sms.SendVerifyCode(phone, verifyCode)
	response.OK(w, map[string]string{
		"ok": "1",
	})
	go models.CreateUserWithPhone(phone)
}

func encodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

type VerifyCodeRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func VerifyUserWithPhone(w http.ResponseWriter, r *http.Request) {
	var request VerifyCodeRequest
	decorder := json.NewDecoder(r.Body)
	err := decorder.Decode(&request)
	if err != nil {
		response.Error(w, http.StatusInternalServerError)
		return
	}
	savedCode, err := cache.Get(request.Phone)
	if err != nil {
		response.Error(w, http.StatusInternalServerError)
		return
	}
	if savedCode != request.Code {
		response.Error(w, http.StatusBadRequest)
		return
	}
	response.OK(w, map[string]string{
		"ok": "1",
	})

}
