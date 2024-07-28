package models

import (
	"time"

	"gorm.io/gorm"
)

type FlowInstance struct {
	ID                  int        `json:"id" gorm:"primaryKey"`
	CreatedAt           time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	StartedAt           *time.Time `json:"started_at"`
	EndedAt             *time.Time `json:"ended_at"`
	ErrorAt             *time.Time `json:"error_at"`
	FlowDefinitionID    *int       `json:"flow_definition_id"`
	FlowDefinitionRefID string     `json:"flow_definition_ref_id"`
	Status              string     `json:"status" gorm:"type:flow_instances_status" default:"not_started"`
	InputData           JSONB      `json:"input_data" gorm:"type:jsonb"`
	OutputData          JSONB      `json:"output_data" gorm:"type:jsonb"`
	Metadata            JSONB      `json:"metadata" gorm:"type:jsonb"`
}

// Constants
const (
	FlowInstanceStatusNotStarted = "not_started"
	FlowInstanceWaiting          = "waiting"
	FlowInstanceStatusRunning    = "running"
	FlowInstanceStatusCompleted  = "completed"
	FlowInstanceStatusFailed     = "failed"
	FlowInstanceStatusStopped    = "stopped"
)

// BeforeCreate hook
func (f *FlowInstance) BeforeCreate(tx *gorm.DB) (err error) {
	// Default status is inactive
	if f.Status == "" {
		f.Status = FlowInstanceStatusNotStarted
	}

	// Default Metadata
	if f.Metadata == nil {
		f.Metadata = make(map[string]interface{})
	}

	// Get FlowDefinition by reference_id the last version
	flowDefinition, err := (&FlowDefinition{}).GetFlowDefinitionByReferenceID(tx, f.FlowDefinitionRefID)
	if err != nil {
		return
	}

	if flowDefinition != nil && f.FlowDefinitionID == nil {
		f.FlowDefinitionID = &flowDefinition.ID
	}

	if f.InputData == nil {
		f.InputData = make(map[string]interface{})
	}

	if f.OutputData == nil {
		f.OutputData = make(map[string]interface{})
	}

	// Default Version
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	f.ErrorAt = nil

	return
}
