package handlers

import (
	"flowspell/models"
	"net/http"
	"strconv"

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

	var flowInstances []models.FlowInstance
	result := h.DB.Limit(limit).Offset(offset).Find(&flowInstances)

	if result.Error != nil {
		h.respondWithError(c, http.StatusInternalServerError, result.Error)
		return
	}
	c.JSON(http.StatusOK, flowInstances)
}

// Create a new flow instance
func (h *FlowInstanceHandler) CreateFlowInstance(c *gin.Context) {
	var flowInstance models.FlowInstance
	if err := c.ShouldBindJSON(&flowInstance); err != nil {
		h.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	// Default status
	if flowInstance.Status == "" {
		flowInstance.Status = models.FlowInstanceStatusNotStarted
	}

	// Default Metadata
	if flowInstance.Metadata == nil {
		flowInstance.Metadata = make(map[string]interface{})
	}

	if result := h.DB.Create(&flowInstance); result.Error != nil {
		h.respondWithError(c, http.StatusInternalServerError, result.Error)
		return
	}
	c.JSON(http.StatusCreated, flowInstance)
}
