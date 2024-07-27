package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
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

// Get FlowDefinition by reference_id the last version
func (f *FlowDefinition) GetFlowDefinitionByReferenceID(db *gorm.DB, referenceID string) (*FlowDefinition, error) {
	var flowDefinition FlowDefinition
	result := db.Where("reference_id = ?", referenceID).Order("version desc").First(&flowDefinition)
	if result.Error != nil {
		return nil, result.Error
	}
	return &flowDefinition, nil
}
