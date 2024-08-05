package models

import (
	"time"
)

type TaskInstance struct {
	ID                  int       `json:"id" gorm:"primaryKey"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Name                string    `json:"name"`
	Status              string    `json:"status" gorm:"type:task_instances_status" default:"not_started"`
	TaskDefinitionID     int      `json:"task_definition_id"`
	TaskDefinitionRefID  string    `json:"task_definition_ref_id"`
	FlowDefinitionID     int       `json:"flow_definition_id"`
	FlowDefinitionRefID  string    `json:"flow_definition_ref_id"`
	ParentTaskID        *int      `json:"parent_task_id"`
	InputData           JSONB     `json:"input_data" gorm:"type:jsonb"`
	OutputData          JSONB     `json:"output_data" gorm:"type:jsonb"`
	Metadata            JSONB     `json:"metadata" gorm:"type:jsonb"`
}
