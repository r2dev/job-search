package auth

import (
	"encoding/json"
	"hirine/cache"
	"net/http"

	"github.com/tj/go/http/response"
)

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
