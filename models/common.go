package models

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"

	"gorm.io/gorm"
)

type SimplifiedSchema struct {
	Properties map[string]map[string]interface{} `json:"properties"`
	Required   []string                          `json:"required"`
}

const flowDefinitionBaseSchema = `{
    "$id": "{{.Host}}/schemas/flow_definitions/{{.ReferenceID}}/{{.Type}}.json",
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "title": "Flow Definition",
    "type": "object",
    "properties": {},
    "required": [],
    "additionalProperties": false
}`

type SchemaData struct {
	Host        string
	ReferenceID string
	Type        string
}

// Constants
const (
    // FlowDefinition
	FlowDefinitionStatusActive   = "active"
	FlowDefinitionStatusInactive = "inactive"

    // FlowInstance
    FlowInstanceStatusWaiting   = "waiting"
	FlowInstanceStatusRunning   = "running"
	FlowInstanceStatusCompleted = "completed"
	FlowInstanceStatusFailed    = "failed"
	FlowInstanceStatusStopped   = "stopped"
)

func CompleteSchema(simplified SimplifiedSchema, schemaData SchemaData) (map[string]interface{}, error) {
	// Replace data
	templateContent, err := template.New("schema").Parse(flowDefinitionBaseSchema)
	if err != nil {
		panic(err)
	}

	var renderedTemplate bytes.Buffer
	err = templateContent.Execute(&renderedTemplate, schemaData)
	if err != nil {
		panic(err)
	}

	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(renderedTemplate.Bytes()), &schema); err != nil {
		return nil, err
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid schema format")
	}

	for key, value := range simplified.Properties {
		properties[key] = value
	}

	schema["required"] = simplified.Required

	return schema, nil
}

// JSONB is a custom type to handle JSONB fields in PostgreSQL
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

// Get FlowDefinition by reference_id the last version
func (f *FlowDefinition) GetFlowDefinitionByReferenceID(db *gorm.DB, referenceID string) (*FlowDefinition, error) {
	var flowDefinition FlowDefinition
	result := db.Where("reference_id = ?", referenceID).Order("version desc").First(&flowDefinition)
	if result.Error != nil {
		return nil, result.Error
	}
	return &flowDefinition, nil
}

// Convert JSONB to String
func ConvertJSONBToString(jsonb JSONB) (string, error) {
	jsonData, err := json.Marshal(jsonb)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// Convert JSONB to map[string]interface{}{}
func ConvertJSONBToMap(jsonb JSONB) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(jsonb)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Convert JSONB to SimplifiedSchema
func ConvertJSONBToSimplifiedSchema(jsonb JSONB) (SimplifiedSchema, error) {
	jsonData, err := json.Marshal(jsonb)
	if err != nil {
		return SimplifiedSchema{}, err
	}
	var simplified SimplifiedSchema
	err = json.Unmarshal(jsonData, &simplified)
	if err != nil {
		return SimplifiedSchema{}, err
	}
	return simplified, nil
}

// Convert SimplifiedSchema to JSONB
func ConvertSimplifiedSchemaToJSONB(simplified SimplifiedSchema, schemaData SchemaData) (JSONB, error) {
	completeSchema, err := CompleteSchema(simplified, schemaData)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(completeSchema)
	if err != nil {
		return nil, err
	}
	var jsonb JSONB
	err = json.Unmarshal(jsonData, &jsonb)
	if err != nil {
		return nil, err
	}
	return jsonb, nil
}
