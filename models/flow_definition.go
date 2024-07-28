package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlowDefinition struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	ReferenceID  string    `json:"reference_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Status       string    `json:"status" gorm:"type:flow_definition_status" default:"inactive"`
	Version      int       `json:"version" default:"1"`
	InputSchema  JSONB     `json:"input_schema" gorm:"type:jsonb"`
	OutputSchema JSONB     `json:"output_schema" gorm:"type:jsonb"`
	Metadata     JSONB     `json:"metadata" gorm:"type:jsonb"`
}

// Constants
const (
	FlowDefinitionStatusActive   = "active"
	FlowDefinitionStatusInactive = "inactive"
)

// BeforeCreate hook
func (f *FlowDefinition) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ReferenceID == "" {
		f.ReferenceID = uuid.New().String()
	}

	// Default status is inactive
	if f.Status == "" {
		f.Status = FlowDefinitionStatusInactive
	}

	// Default Metadata
	if f.Metadata == nil {
		f.Metadata = make(map[string]interface{})
	}

	if f.InputSchema == nil {
		f.InputSchema = make(map[string]interface{})
	}

	if f.OutputSchema == nil {
		f.OutputSchema = make(map[string]interface{})
	}

	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	f.Version = 1

	return
}
