package stdext

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// Use json to convert the interface{} to a map if possible. Return json errors if any.
func ToMap(s interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type SqlJson map[string]interface{}

func (s *SqlJson) Scan(src interface{}) error {
	val, ok := src.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("could not scan %v as JSON", src))
	}
	var result map[string]interface{}
	err := json.Unmarshal(val, &result)
	*s = result
	return err
}

func (s SqlJson) Value() (driver.Value, error) {
	return json.Marshal(s)
}
