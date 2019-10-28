package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
)

var decodeJSONError = errors.New("request is not able to decode")

func DecodeJSON(r *http.Request, result interface{}) error {
	decorder := json.NewDecoder(r.Body)
	err := decorder.Decode(&result)
	if err != nil {
		return decodeJSONError
	}
	return nil
}
