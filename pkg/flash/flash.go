package flash

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
)

type FlashMessageType int

const (
	Success = iota + 1
	Fail
	Warning
)

type FlashMessage struct {
	Type    FlashMessageType
	Message string
}

var FlashEncodeError = errors.New("flash cant encode")
var FlashDecodeError = errors.New("flash decode error")

func Encode(fm FlashMessage) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(fm)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func Decode(str string) (FlashMessage, error) {
	m := FlashMessage{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return FlashMessage{}, FlashDecodeError
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		return FlashMessage{}, FlashDecodeError
	}
	return m, nil
}
