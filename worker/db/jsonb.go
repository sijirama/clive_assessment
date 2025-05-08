package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONBMap is a custom type for handling JSONB in PostgreSQL with map[string]interface{}
type JSONBMap map[string]interface{}

// Value implements the driver.Valuer interface to convert JSONBMap to a database value
func (j JSONBMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface to convert database value to JSONBMap
func (j *JSONBMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("scan source is not []byte")
	}
	if len(bytes) == 0 {
		*j = make(map[string]interface{})
		return nil
	}
	return json.Unmarshal(bytes, j)
}
