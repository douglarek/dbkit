package dbc

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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

/************************************/
/****** From sqlx types *************/
/************************************/

// JSONText is a json.RawMessage, which is a []byte underneath.
// Value() validates the json format in the source, and returns an error if
// the json is not valid.  Scan does no validation.  JSONText additionally
// implements `Unmarshal`, which unmarshals the json within to an interface{}
type JSONText json.RawMessage

// MarshalJSON returns j as the JSON encoding of j.
func (j JSONText) MarshalJSON() ([]byte, error) {
	return j, nil
}

// UnmarshalJSON sets *j to a copy of data
func (j *JSONText) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("JSONText: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil

}

// Value returns j as a value.  This does a validating unmarshal into another
// RawMessage.  If j is invalid json, it returns an error.
func (j JSONText) Value() (driver.Value, error) {
	var m json.RawMessage
	var err = j.Unmarshal(&m)
	if err != nil {
		return []byte{}, err
	}
	return string(j), nil // modified here: returns string, not bytes!
}

// Scan stores the src in *j.  No validation is done.
func (j *JSONText) Scan(src interface{}) error {
	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return errors.New("Incompatible type for JSONText")
	}
	*j = JSONText(append((*j)[0:0], source...))
	return nil
}

// Unmarshal unmarshal's the json in j to v, as in json.Unmarshal.
func (j *JSONText) Unmarshal(v interface{}) error {
	return json.Unmarshal([]byte(*j), v)
}

// Pretty printing for JSONText types
func (j JSONText) String() string {
	return string(j)
}
