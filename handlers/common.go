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

type FlowDefinitionRequestBody struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Input       models.JSONB `json:"input"`
	Output      models.JSONB `json:"output" gorm:"-"`
}
