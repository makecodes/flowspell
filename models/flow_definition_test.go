package models

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestFlowDefinition(t *testing.T) {
    fd := FlowDefinition{
        Name: "Test Flow",
    }
    assert.Equal(t, "Test Flow", fd.Name)
}
