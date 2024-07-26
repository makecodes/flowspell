package models

import (
    "time"
)

type FlowDefinition struct {
    ID          int       `json:"id" gorm:"primaryKey"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    Status      string    `json:"status" gorm:"type:flow_definition_status"`
    Metadata    JSONB     `json:"metadata" gorm:"type:jsonb"`
}
