package handlers

import (
	"errors"
	"flowspell/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FlowDefinitionHandler struct {
	DB *gorm.DB
}

// Helper function to respond with an error
func (h *FlowDefinitionHandler) respondWithError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}

// GetFlowDefinitions Return all flow definitions
func (h *FlowDefinitionHandler) GetFlowDefinitions(c *gin.Context) {
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

	var flowDefinitions []models.FlowDefinition
	result := h.DB.Order("created_at DESC").Limit(limit).Offset(offset).Find(&flowDefinitions)

	if result.Error != nil {
		h.respondWithError(c, http.StatusInternalServerError, result.Error)
		return
	}

	response := make([]FlowDefinitionResponse, len(flowDefinitions))
	for i, fd := range flowDefinitions {
		response[i] = serializeFlowDefinition(fd)
	}

	c.JSON(http.StatusOK, response)
}

// CreateFlowDefinition Create a new flow definition
func (h *FlowDefinitionHandler) CreateFlowDefinition(c *gin.Context) {
	var flowDefinition models.FlowDefinition
	if err := c.ShouldBindJSON(&flowDefinition); err != nil {
		h.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	// Verify if the definition name already exists
	var count int64
	h.DB.Model(&models.FlowDefinition{}).Where("name = ?", flowDefinition.Name).Count(&count)
	if count > 0 {
		err := &CustomError{
			Message: map[string]string{
				"name": "A flow with this name already exists",
			},
		}
		h.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	if result := h.DB.Create(&flowDefinition); result.Error != nil {
		h.respondWithError(c, http.StatusInternalServerError, result.Error)
		return
	}

	response := serializeFlowDefinition(flowDefinition)

	c.JSON(http.StatusCreated, response)
}

// UpdateFlowDefinition Update a flow definition by its ID
func (h *FlowDefinitionHandler) UpdateFlowDefinition(c *gin.Context) {
	referenceId := c.Param("referenceId")
	flowDefinition, err := models.GetLastFlowDefinitionVersionFromReferenceID(h.DB, referenceId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.respondWithError(c, http.StatusNotFound, err)
		} else {
			h.respondWithError(c, http.StatusInternalServerError, err)
		}
		return
	}

	if err := c.ShouldBindJSON(&flowDefinition); err != nil {
		h.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	flowDefinition.ID = 0

	if result := h.DB.Create(&flowDefinition); result.Error != nil {
		h.respondWithError(c, http.StatusInternalServerError, result.Error)
		return
	}

	response := serializeFlowDefinition(flowDefinition)

	c.JSON(http.StatusCreated, response)
}

// GetFlowDefinition Get a flow definition by its ID
func (h *FlowDefinitionHandler) GetFlowDefinition(c *gin.Context) {
	referenceId := c.Param("referenceId")
	flowDefinition, err := models.GetLastFlowDefinitionVersionFromReferenceID(h.DB, referenceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.respondWithError(c, http.StatusNotFound, err)
		} else {
			h.respondWithError(c, http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, flowDefinition)
}

// DeleteFlowDefinition Delete a flow definition by its ID
func (h *FlowDefinitionHandler) DeleteFlowDefinition(c *gin.Context) {
	referenceId := c.Param("referenceId")
	result := h.DB.Where("reference_id = ?", referenceId).Delete(&models.FlowDefinition{})
	if result.Error != nil {
		h.respondWithError(c, http.StatusInternalServerError, result.Error)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetFlowDefinitionSchema handle the /schemas/flow_definitions/:referenceId/:type.json endpoint
func (h *FlowDefinitionHandler) GetFlowDefinitionSchema(c *gin.Context) {
	referenceId := c.Param("referenceId")
	flowDefinition, err := models.GetLastFlowDefinitionVersionFromReferenceID(h.DB, referenceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.respondWithError(c, http.StatusNotFound, err)
		} else {
			h.respondWithError(c, http.StatusInternalServerError, err)
		}
		return
	}

	var schema models.JSONB
	schemaType := c.Param("type")

	switch schemaType {
	case "input.json":
		schema = flowDefinition.InputSchema
	case "output.json":
		schema = flowDefinition.OutputSchema
	default:
		h.respondWithError(c, http.StatusBadRequest, &CustomError{
			Message: map[string]string{
				"type": "Invalid schema type",
			},
		})
		return
	}

	c.JSON(http.StatusOK, schema)
}
