package models

import (
	"errors"
	"time"
)

type CustomTime struct {
	time.Time
}

const customLayout = "2006-01-02T15:04:05.999999" // format tanpa timezone

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		return nil
	}

	// remove quotes
	s = s[1 : len(s)-1]

	t, err := time.Parse(customLayout, s)
	if err != nil {
		return errors.New("invalid time format: " + err.Error())
	}

	ct.Time = t
	return nil
}
