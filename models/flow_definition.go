package models

import (
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlowDefinition struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	ReferenceID  string    `json:"reference_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
    IsLatest     bool      `json:"is_latest"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Status       string    `json:"status" gorm:"type:flow_definition_status" default:"inactive"`
	Version      int       `json:"version" default:"1"`
	Input        JSONB     `json:"input,omitempty" gorm:"-"`
	Output       JSONB     `json:"output,omitempty" gorm:"-"`
	InputSchema  JSONB     `json:"input_schema" gorm:"type:jsonb"`
	OutputSchema JSONB     `json:"output_schema" gorm:"type:jsonb"`
	Metadata     JSONB     `json:"metadata" gorm:"type:jsonb"`
}

// BeforeCreate hook
func (f *FlowDefinition) BeforeCreate(tx *gorm.DB) (err error) {
	flowspellHost := os.Getenv("FLOWSPELL_HOST")
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

	// Verify if the f.Input schema is valid
	convertedInput, err := ConvertJSONBToSimplifiedSchema(f.Input)
	if err != nil {
		return err
	}

	schemaDataInput := SchemaData{
		Host:        flowspellHost,
		ReferenceID: f.ReferenceID,
		Type:        "input",
	}

	completeInputSchema, err := CompleteSchema(convertedInput, schemaDataInput)
	if err != nil {
		return err
	}

	f.InputSchema = completeInputSchema

	// Verify if the f.Output schema is valid
	convertedOutput, err := ConvertJSONBToSimplifiedSchema(f.Output)
	if err != nil {
		return err
	}

	schemaDataOutput := SchemaData{
		Host:        flowspellHost,
		ReferenceID: f.ReferenceID,
		Type:        "output",
	}

	completeOutputSchema, err := CompleteSchema(convertedOutput, schemaDataOutput)
	if err != nil {
		return err
	}

	f.OutputSchema = completeOutputSchema

	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	if f.Version == 0 {
		f.Version = 1
	}

	f.Version++

    f.IsLatest = true

    // Set all other versions to false
    tx.Model(&FlowDefinition{}).Where("reference_id = ?", f.ReferenceID).Update("is_latest", false)

	return
}

func (f *FlowDefinition) CountTaskDefinitionsByFlowDefinitionRefID(tx *gorm.DB) (count int64, err error) {
    err = tx.
        Model(&TaskDefinition{}).
        Where("flow_definition_ref_id = ?", f.ReferenceID).
        Where("is_latest = ?", true).
        Count(&count).Error

    return
}

func GetLastFlowDefinitionVersionFromReferenceID(tx *gorm.DB, referenceID string) (flowDefinition FlowDefinition, err error) {
    err = tx.
        Where("reference_id = ?", referenceID).
        Order("version desc").
        First(&flowDefinition).
        Error

    return
}
