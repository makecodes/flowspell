package models

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
)

// JSONB is a custom type to handle JSONB fields in PostgreSQL
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
    return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("type assertion to []byte failed")
    }
    return json.Unmarshal(bytes, j)
}
