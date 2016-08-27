package dbc

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

// TimeStamp is a unix-like timestamp.
type TimeStamp struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface.
func (t TimeStamp) MarshalJSON() ([]byte, error) {
	var s string
	if t.Time.IsZero() {
		s = `0` // maybe others ? up to u :)
	} else {
		s = fmt.Sprintf("%v", t.Time.Unix())
	}
	return []byte(s), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *TimeStamp) UnmarshalJSON(data []byte) (err error) {
	ts, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return
	}
	t.Time = time.Unix(ts, 0)
	return
}

// Scan implements the Scanner interface.
func (t *TimeStamp) Scan(value interface{}) error {
	if value != nil {
		t.Time = value.(time.Time)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (t TimeStamp) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}
