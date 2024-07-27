package models

import (
	"time"
)

type FlowInstance struct {
	ID               int            `json:"id" gorm:"primaryKey"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	StartedAt        time.Time      `json:"started_at"`
	FlowDefinitionID int            `json:"flow_definition_id"`
	FlowDefinition   FlowDefinition `json:"flow_definition" gorm:"foreignKey:FlowDefinitionID"`
	Status           string         `json:"status" gorm:"type:flow_instances_status" default:"not_started"`
	Version          int            `json:"version" default:"1"`
	Metadata         JSONB          `json:"metadata" gorm:"type:jsonb"`
}

// Constants
const (
	FlowInstanceStatusNotStarted = "not_started"
	FlowInstanceStatusRunning    = "running"
	FlowInstanceStatusCompleted  = "completed"
	FlowInstanceStatusFailed     = "failed"
	FlowInstanceStatusStopped    = "stopped"
)
