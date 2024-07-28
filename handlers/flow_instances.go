package handlers

import (
	"flowspell/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FlowInstanceHandler struct {
	DB *gorm.DB
}

// Helper function to respond with an error
func (h *FlowInstanceHandler) respondWithError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}

// Find a flow instance by its ID
func (h *FlowInstanceHandler) findFlowDefinitionByID(id string) (*models.FlowInstance, error) {
	var flowInstance models.FlowInstance
	if result := h.DB.First(&flowInstance, id); result.Error != nil {
		return nil, result.Error
	}
	return &flowInstance, nil
}

// Return all flow instances
func (h *FlowInstanceHandler) GetFlowInstances(c *gin.Context) {
	limit := 25
	offset := 0
	flowDefinitionID := 0

	if c.Query("limit") != "" {
		if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
			limit = l
		}
	}

	if limit > 1000 {
		limit = 1000
	}

	if c.Query("offset") != "" {
		if o, err := strconv.Atoi(c.Query("offset")); err == nil && o >= 0 {
			offset = o
		}
	}

	if c.Query("flow_definition") != "" {
		_, err := strconv.Atoi(c.Query("flow_definition"))
		if err != nil {
			h.respondWithError(c, http.StatusBadRequest, err)
			return
		}
	}

	var flowInstances []models.FlowInstance
	result := h.DB.Limit(limit).Offset(offset).Find(&flowInstances)

	if flowDefinitionID != 0 {
		result = h.DB.Where("flow_definition_id = ?", flowDefinitionID).Limit(limit).Offset(offset).Find(&flowInstances)
	}

	if result.Error != nil {
		h.respondWithError(c, http.StatusInternalServerError, result.Error)
		return
	}
	c.JSON(http.StatusOK, flowInstances)
}

// // Create a new flow instance
// func (h *FlowInstanceHandler) CreateFlowInstance(c *gin.Context) {
// 	var flowInstance models.FlowInstance
// 	if err := c.ShouldBindJSON(&flowInstance); err != nil {
// 		h.respondWithError(c, http.StatusBadRequest, err)
// 		return
// 	}

// 	if result := h.DB.Create(&flowInstance); result.Error != nil {
// 		h.respondWithError(c, http.StatusInternalServerError, result.Error)
// 		return
// 	}
// 	c.JSON(http.StatusCreated, flowInstance)
// }

// Start a flow instance
func (h *FlowInstanceHandler) StartFlowInstance(c *gin.Context) {
	id := c.Param("referenceId")
	flowInstance, err := h.findFlowDefinitionByID(id)
	if err != nil {
		h.respondWithError(c, http.StatusNotFound, err)
		return
	}

	now := time.Now()
	flowInstance.StartedAt = &now
	flowInstance.Status = models.FlowInstanceStatusRunning
}
