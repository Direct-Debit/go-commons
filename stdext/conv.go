package stdext

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// ToMap uses json to convert the interface{} to a map if possible. Return json errors if any.
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

type SqlJsonArray []interface{}

func (s *SqlJsonArray) Scan(src interface{}) error {
	val, ok := src.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("could not scan %v as JSON", src))
	}
	var result []interface{}
	err := json.Unmarshal(val, &result)
	*s = result
	return err
}

func (s SqlJsonArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func SliceInterfaceToString(list []interface{}) []string {
	result := make([]string, len(list))
	for i, l := range list {
		result[i] = l.(string)
	}
	return result
}

func Ptr[T any](val T) *T {
	return &val
}

func DerefSafe[T any](t *T) T {
	if t == nil || t == (*T)(nil) {
		var nothing T
		return nothing
	}
	return *t
}
