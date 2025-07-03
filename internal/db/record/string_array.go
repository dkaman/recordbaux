package record

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type StringArray []string

func (sa StringArray) Value() (driver.Value, error) {
	if sa == nil {
		return nil, nil
	}
	j, err := json.Marshal(sa)
	return driver.Value(j), err
}

func (sa *StringArray) Scan(src interface{}) error {
	if src == nil {
		*sa = StringArray{}
		return nil
	}
	var source []byte
	switch src := src.(type) {
	case string:
		source = []byte(src)
	case []byte:
		source = src
	default:
		return errors.New("incompatible type for StringArray")
	}
	return json.Unmarshal(source, sa)
}
