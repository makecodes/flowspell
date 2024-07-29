package handlers

import (
	"flowspell/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskDefinitionHandler struct {
	DB *gorm.DB
}

// Helper function to respond with an error
func (h *TaskDefinitionHandler) respondWithError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}


// Return all task definitions
func (h *TaskDefinitionHandler) GetTaskDefinitions(c *gin.Context) {
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

	var taskDefinitions []models.TaskDefinition
	result := h.DB.Order("created_at DESC").Limit(limit).Offset(offset).Find(&taskDefinitions)

	if result.Error != nil {
		h.respondWithError(c, http.StatusInternalServerError, result.Error)
		return
	}
	c.JSON(http.StatusOK, taskDefinitions)
}


// Create a new task definition
func (h *TaskDefinitionHandler) CreateTaskDefinition(c *gin.Context) {
	var taskDefinition models.TaskDefinition
	if err := c.ShouldBindJSON(&taskDefinition); err != nil {
		h.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	// Verify if the definition name already exists
	var count int64
	h.DB.Model(&models.TaskDefinition{}).Where("name = ?", taskDefinition.Name).Count(&count)
	if count > 0 {
		err := &CustomError{
			Message: map[string]string{
				"name": "A task with this name already exists",
			},
		}
		h.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	if result := h.DB.Create(&taskDefinition); result.Error != nil {
		h.respondWithError(c, http.StatusInternalServerError, result.Error)
		return
	}
	response := TaskDefinitionResponse{
		ID:              taskDefinition.ID,
		CreatedAt:       taskDefinition.CreatedAt,
		UpdatedAt:       taskDefinition.UpdatedAt,
		ReferenceID:     taskDefinition.ReferenceID,
		Name:            taskDefinition.Name,
		Description:     taskDefinition.Description,
        FlowDefinitionRefID: taskDefinition.FlowDefinitionRefID,
        ParentTaskID:    *taskDefinition.ParentTaskID,
		InputSchema:     taskDefinition.InputSchema,
		OutputSchema:    taskDefinition.OutputSchema,
		Version:         taskDefinition.Version,
		Metadata:        taskDefinition.Metadata,
	}

	c.JSON(http.StatusCreated, response)
}
