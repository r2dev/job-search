package helpers

import "time"

func ParseJavascriptTimeString(str string) (t time.Time) {
	layout := "2006-01-02T15:04:05.000Z"

	t, err := time.Parse(layout, str)
	if err != nil {
		t = time.Time{}
	}
	return
}
