package models

import (
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskDefinition struct {
	ID                  int       `json:"id" gorm:"primaryKey"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    IsLatest            bool      `json:"is_latest"`
	ReferenceID         string    `json:"reference_id"`
	Name                string    `json:"name"`
	Description         string    `json:"description,omitempty"`
	FlowDefinitionID     int       `json:"flow_definition_id"`
	FlowDefinitionRefID  string    `json:"flow_definition_ref_id"`
	ParentTaskID        *int      `json:"parent_task_id"`
    ParentTaskDefRefID  *string   `json:"parent_task_def_ref_id"`
	InputSchema         JSONB     `json:"input_schema" gorm:"type:jsonb"`
	OutputSchema        JSONB     `json:"output_schema" gorm:"type:jsonb"`
	Input               JSONB     `json:"input,omitempty" gorm:"-"`
	Output              JSONB     `json:"output,omitempty" gorm:"-"`
	Version             int       `json:"version" default:"1"`
	Metadata            JSONB     `json:"metadata" gorm:"type:jsonb"`
}

// BeforeCreate hook
func (f *TaskDefinition) BeforeCreate(tx *gorm.DB) (err error) {
	flowspellHost := os.Getenv("FLOWSPELL_HOST")
	if f.ReferenceID == "" {
		f.ReferenceID = uuid.New().String()
	}

	// Get FlowDefinition by reference_id the last version
	flowDefinition, err := (&FlowDefinition{}).GetFlowDefinitionByReferenceID(tx, f.FlowDefinitionRefID)
	if err != nil {
		return
	}

	f.FlowDefinitionID = flowDefinition.ID

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
		Type:        "input",
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
    tx.Model(&TaskDefinition{}).
        Where("reference_id = ?", f.ReferenceID).
        Update("is_latest", false)

	return
}


// func (f *TaskDefinition) TaskDefinitionsByFlowDefinitionRefID(tx *gorm.DB, referenceID string) (taskDefinitions []TaskDefinition, err error) {
//     result := tx.Where("flow_definition_ref_id = ?", referenceID).Find(&taskDefinitions)
//     if result.Error != nil {
//         return nil, result.Error
//     }
//     return taskDefinitions, nil
// }


func GetLastTaskDefinitionVersionFromReferenceID(tx *gorm.DB, referenceID string) (taskDefinition TaskDefinition, err error) {
    err = tx.Where("reference_id = ?", referenceID).Order("version desc").First(&taskDefinition).Error

    return
}
