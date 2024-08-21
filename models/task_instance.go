package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskInstance struct {
	ID                  int        `json:"id" gorm:"primaryKey"`
	CreatedAt           time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	AcknowledgedAt      *time.Time `json:"acknowledged_at"`
	StartedAt           *time.Time `json:"started_at"`
	CompletedAt         *time.Time `json:"completed_at"`
	FailedAt            *time.Time `json:"failed_at"`
	Name                string     `json:"name"`
	Status              string     `json:"status" gorm:"type:task_instances_status" default:"not_started"`
	TaskDefinitionID    int        `json:"task_definition_id"`
	TaskDefinitionRefID string     `json:"task_definition_ref_id"`
	FlowDefinitionID    int        `json:"flow_definition_id"`
	FlowDefinitionRefID string     `json:"flow_definition_ref_id"`
	ParentTaskID        *int       `json:"parent_task_id"`
	InputData           JSONB      `json:"input_data" gorm:"type:jsonb"`
	OutputData          JSONB      `json:"output_data" gorm:"type:jsonb"`
	Metadata            JSONB      `json:"metadata" gorm:"type:jsonb"`
}


func GetAcknowledgedTasks(tx *gorm.DB) ([]TaskInstance, error) {
    var taskInstances []TaskInstance
    err := tx.
        Where("status = ? AND acknowledged_at < ?", TaskInstanceAcknowledged, time.Now().
        Add(-5*time.Minute)).
        Find(&taskInstances).
        Error
    return taskInstances, err
}
