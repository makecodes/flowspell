package models

import (
    "time"
)

type FlowInstance struct {
    ID               int            `json:"id" gorm:"primaryKey"`
    CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
    FlowDefinitionID int            `json:"flow_definition_id"`
    FlowDefinition   FlowDefinition `json:"flow_definition" gorm:"foreignKey:FlowDefinitionID"`
    Status           string         `json:"status" gorm:"type:flow_instances_status"`
    Metadata         JSONB          `json:"metadata" gorm:"type:jsonb"`
}

