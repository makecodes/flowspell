package models

import (
	"time"
)

type TaskDefinition struct {
	ID           int             `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Type         string          `json:"type" gorm:"type:task_definition_types"`
	ParentTaskID int             `json:"parent_task_id"`
	ParentTask   *TaskDefinition `json:"parent_task" gorm:"foreignKey:ParentTaskID"`
	InputSchema  JSONB           `json:"input_schema" gorm:"type:jsonb"`
	OutputSchema JSONB           `json:"output_schema" gorm:"type:jsonb"`
}
