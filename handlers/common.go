package handlers

import (
	"flowspell/models"
	"fmt"
	"time"
)

type CustomError struct {
	Message map[string]string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

type FlowDefinitionResponse struct {
	ID           int          `json:"id"`
	ReferenceID  string       `json:"reference_id"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Status       string       `json:"status"`
	Version      int          `json:"version"`
	InputSchema  models.JSONB `json:"input_schema"`
	OutputSchema models.JSONB `json:"output_schema"`
	Metadata     models.JSONB `json:"metadata"`
}

type TaskDefinitionResponse struct {
	ID                  int          `json:"id"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
	ReferenceID         string       `json:"reference_id"`
	Name                string       `json:"name"`
	Description         string       `json:"description"`
	FlowDefinitionRefID  string       `json:"flow_definition_ref_id"`
	ParentTaskID        int          `json:"parent_task_id"`
	InputSchema         models.JSONB `json:"input_schema"`
	OutputSchema        models.JSONB `json:"output_schema"`
	Version             int          `json:"version"`
	Metadata            models.JSONB `json:"metadata"`
}

type TaskDefinitionRequestBody struct {
	Name                string       `json:"name"`
	Description         *string      `json:"description"`
	FlowDefinitionRefID string       `json:"flow_definition_ref_id"`
	ParentTaskID        *int         `json:"parent_task_id"`
	Input               models.JSONB `json:"input"`
	Output              models.JSONB `json:"output" gorm:"-"`
}

type FlowDefinitionRequestBody struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Input       models.JSONB `json:"input"`
	Output      models.JSONB `json:"output" gorm:"-"`
}

// Find a task definition by its ID
func (h *TaskDefinitionHandler) findTaskDefinitionByReferenceID(refrenceId string) (*models.TaskDefinition, error) {
	var taskDefinition models.TaskDefinition
	result := h.DB.Where("reference_id = ?", refrenceId).Order("version desc").First(&taskDefinition)
	if result.Error != nil {
		return nil, result.Error
	}
	return &taskDefinition, nil
}
