package models

import (
	"database/sql/driver"
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

// MarshalJSON ensures CustomTime is marshaled in the same layout
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.IsZero() {
		return []byte("null"), nil
	}
	return []byte("\"" + ct.Format(customLayout) + "\""), nil
}

// Value implements the driver.Valuer interface so CustomTime can be stored in DB
func (ct CustomTime) Value() (driver.Value, error) {
	if ct.IsZero() {
		return nil, nil
	}
	return ct.Time, nil
}

// Scan implements the sql.Scanner interface so CustomTime can be read from DB
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		ct.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		ct.Time = v
		return nil
	case []byte:
		s := string(v)
		// try custom layout first, then RFC3339
		if t, err := time.Parse(customLayout, s); err == nil {
			ct.Time = t
			return nil
		}
		if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
			ct.Time = t
			return nil
		}
		return errors.New("cannot parse time from bytes")
	case string:
		s := v
		if t, err := time.Parse(customLayout, s); err == nil {
			ct.Time = t
			return nil
		}
		if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
			ct.Time = t
			return nil
		}
		return errors.New("cannot parse time from string")
	default:
		return errors.New("unsupported Scan value for CustomTime")
	}
}

// GormDataType allows GORM to recognize the data type for migrations
func (CustomTime) GormDataType() string {
	return "timestamp"
}
