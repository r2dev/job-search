package auth

import (
	"crypto/rand"
	"encoding/json"
	"hirine/cache"
	"hirine/models"
	"hirine/sms"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tj/go/http/response"
)

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
		log.WithError(err)
		response.Error(w, http.StatusInternalServerError)
		return
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
