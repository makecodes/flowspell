package models

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestFlowInstance(t *testing.T) {
    fi := FlowInstance{
        ID: 1,
    }
    assert.Equal(t, 1, fi.ID)
}
