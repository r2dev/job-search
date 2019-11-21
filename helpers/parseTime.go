package helpers

import "time"

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
