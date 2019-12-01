package helpers

import (
	"strconv"
	"time"
)

func ParseJavascriptTimeString(str string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"

	t, err := time.Parse(layout, str)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func ParseDateTimeLocalString(str string) (time.Time, error) {
	layout := "2006-01-02T15:04"
	t, err := time.Parse(layout, str)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func ParseUnixString(str string) (time.Time, error) {
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	t := time.Unix(v, 0)
	return t, nil
}
