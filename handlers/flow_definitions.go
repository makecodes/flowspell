package handlers

import (
    "flowspell/models"
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type FlowDefinitionHandler struct {
    DB *gorm.DB
}

// Find a flow definition by its ID
func (h *FlowDefinitionHandler) findFlowDefinitionByID(id string) (*models.FlowDefinition, error) {
    var flowDefinition models.FlowDefinition
    if result := h.DB.First(&flowDefinition, id); result.Error != nil {
        return nil, result.Error
    }
    return &flowDefinition, nil
}

// Helper function to respond with an error
func (h *FlowDefinitionHandler) respondWithError(c *gin.Context, status int, err error) {
    c.JSON(status, gin.H{"error": err.Error()})
}

// Return all flow definitions
func (h *FlowDefinitionHandler) GetFlowDefinitions(c *gin.Context) {
    var flowDefinitions []models.FlowDefinition
    if result := h.DB.Find(&flowDefinitions); result.Error != nil {
        h.respondWithError(c, http.StatusInternalServerError, result.Error)
        return
    }
    c.JSON(http.StatusOK, flowDefinitions)
}

// Create a new flow definition
func (h *FlowDefinitionHandler) CreateFlowDefinition(c *gin.Context) {
    var flowDefinition models.FlowDefinition
    if err := c.ShouldBindJSON(&flowDefinition); err != nil {
        h.respondWithError(c, http.StatusBadRequest, err)
        return
    }
    if result := h.DB.Create(&flowDefinition); result.Error != nil {
        h.respondWithError(c, http.StatusInternalServerError, result.Error)
        return
    }
    c.JSON(http.StatusCreated, flowDefinition)
}

// Get a flow definition by its ID
func (h *FlowDefinitionHandler) GetFlowDefinition(c *gin.Context) {
    id := c.Param("id")
    flowDefinition, err := h.findFlowDefinitionByID(id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            h.respondWithError(c, http.StatusNotFound, err)
        } else {
            h.respondWithError(c, http.StatusInternalServerError, err)
        }
        return
    }
    c.JSON(http.StatusOK, flowDefinition)
}

// Delete a flow definition by its ID
func (h *FlowDefinitionHandler) DeleteFlowDefinition(c *gin.Context) {
    id := c.Param("id")
    if result := h.DB.Delete(&models.FlowDefinition{}, id); result.Error != nil {
        h.respondWithError(c, http.StatusInternalServerError, result.Error)
        return
    }
    c.JSON(http.StatusNoContent, nil)
}

// Update a flow definition by its ID
func (h *FlowDefinitionHandler) UpdateFlowDefinition(c *gin.Context) {
    id := c.Param("id")
    flowDefinition, err := h.findFlowDefinitionByID(id)
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

    if result := h.DB.Save(&flowDefinition); result.Error != nil {
        h.respondWithError(c, http.StatusInternalServerError, result.Error)
        return
    }

    c.JSON(http.StatusOK, flowDefinition)
}
