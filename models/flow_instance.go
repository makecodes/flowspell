package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/kaptinlin/jsonschema"
	"gorm.io/gorm"
)

type FlowInstance struct {
	ID                  int        `json:"id" gorm:"primaryKey"`
	CreatedAt           time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	StartedAt           *time.Time `json:"started_at"`
	EndedAt             *time.Time `json:"ended_at"`
	ErrorAt             *time.Time `json:"error_at"`
	FlowDefinitionID     *int       `json:"flow_definition_id"`
	FlowDefinitionRefID  string     `json:"flow_definition_ref_id"`
	Status              string     `json:"status" gorm:"type:flow_instances_status" default:"waiting"`
	InputData           JSONB      `json:"input_data" gorm:"type:jsonb"`
	OutputData          JSONB      `json:"output_data" gorm:"type:jsonb"`
	Metadata            JSONB      `json:"metadata" gorm:"type:jsonb"`
}

// BeforeCreate hook
func (f *FlowInstance) BeforeCreate(tx *gorm.DB) (err error) {
	// Default status is inactive
	if f.Status == "" {
		f.Status = FlowInstanceStatusWaiting
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

    if flowDefinition.Status == FlowDefinitionStatusInactive {
        err = errors.New("flow definition is not active")
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

	// Validating the input data against the input schema
	inputJSONSchema, err := ConvertJSONBToString(flowDefinition.InputSchema)
	if err != nil {
		return err
	}

	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile([]byte(inputJSONSchema))
	if err != nil {
		return err
	}

	inputData, err := ConvertJSONBToMap(f.InputData)
	if err != nil {
		return
	}

	result := schema.Validate(inputData)
	if !result.IsValid() {
		schemaUrl := flowDefinition.InputSchema["$id"]
		err = errors.New("input data is not valid, refer to the schema: " + fmt.Sprintf("%v", schemaUrl))
		return err
	}

	if f.Status == FlowInstanceStatusCompleted {
		// Validating the output data against the output schema
		outputJSONSchema, err := ConvertJSONBToString(flowDefinition.OutputSchema)
		if err != nil {
			return err
		}

		compiler = jsonschema.NewCompiler()
		schema, err = compiler.Compile([]byte(outputJSONSchema))
		if err != nil {
			return err
		}

		outputData, _ := ConvertJSONBToMap(f.OutputData)
		result := schema.Validate(outputData)
		if !result.IsValid() {
			schemaUrl := flowDefinition.OutputSchema["$id"]
			err = errors.New("input data is not valid, refer to the schema: " + fmt.Sprintf("%v", schemaUrl))
			return err
		}
	}

	// Default Version
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	f.ErrorAt = nil

	return
}
