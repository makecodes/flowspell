package handlers

import (
	"encoding/json"
	"flowspell/models"
	"io"
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

// Start a flow instance
func (h *FlowInstanceHandler) StartFlow(c *gin.Context) {
	referenceId := c.Param("referenceId")

	// Get the flow definition by reference ID
	var flowDefinition models.FlowDefinition
	flowDefinition, err := models.GetLastFlowDefinitionVersionFromReferenceID(h.DB, referenceId)
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	// Verify if the flow definition is active
	if flowDefinition.Status != models.FlowDefinitionStatusActive {
		err := &CustomError{
			Message: map[string]string{
				"status": "Flow definition is not active",
			},
		}
		h.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	// Retrieve all task definitions by flow definition reference ID
	taskDefinitions, err := models.GetTaskDefinitionsByFlowDefinitionRefID(h.DB, referenceId)
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	// Verify if the flow definition has tasks
	if len(taskDefinitions) == 0 {
		err := &CustomError{
			Message: map[string]string{
				"name": "Flow definition has no tasks",
			},
		}
		h.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	go func() {

		// Create task instances with status not_started
		for _, taskDefinition := range taskDefinitions {
			var taskInstance models.TaskInstance
			taskInstance.Name = taskDefinition.Name
			taskInstance.TaskDefinitionID = taskDefinition.ID
			taskInstance.TaskDefinitionRefID = taskDefinition.ReferenceID
			taskInstance.FlowDefinitionID = flowDefinition.ID
			taskInstance.FlowDefinitionRefID = flowDefinition.ReferenceID
			taskInstance.Status = models.TaskInstanceStatusNotStarted

			if taskDefinition.ParentTaskID != nil {
				taskInstance.ParentTaskID = taskDefinition.ParentTaskID
			}

			if result := h.DB.Create(&taskInstance); result.Error != nil {
				h.respondWithError(c, http.StatusInternalServerError, result.Error)
				return
			}
		}
	}()

	// Create a new flow instance
	var flowInstance models.FlowInstance
	if err := c.ShouldBindJSON(&flowInstance); err != nil {
		h.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	flowInstance.FlowDefinitionRefID = referenceId

	now := time.Now()

	flowInstance.StartedAt = &now
	flowInstance.Status = models.FlowInstanceStatusRunning

	if result := h.DB.Create(&flowInstance); result.Error != nil {
		h.respondWithError(c, http.StatusInternalServerError, result.Error)
		return
	}
	c.JSON(http.StatusCreated, flowInstance)
}

func GetJSONFromRequestBody(body io.ReadCloser) (map[string]interface{}, error) {
	var inputData map[string]interface{}
	if err := json.NewDecoder(body).Decode(&inputData); err != nil {
		return nil, err
	}
	return inputData, nil
}
