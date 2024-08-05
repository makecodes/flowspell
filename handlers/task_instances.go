package handlers

import (
	"flowspell/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskInstanceHandler struct {
	DB *gorm.DB
}

// Helper function to respond with an error
func (h *TaskInstanceHandler) respondWithError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}

// Return all task instances
func (h *TaskInstanceHandler) GetTaskQueue(c *gin.Context) {
	limit := 100
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

	var taskInstances []models.TaskInstance
	result := h.DB.Where("Status = ?", models.TaskInstanceStatusNotStarted).Limit(limit).Offset(offset).Find(&taskInstances)

	if result.Error != nil {
		h.respondWithError(c, http.StatusInternalServerError, result.Error)
		return
	}

	// use a goroutine to update the task instances and mark them as acknowledged
	go func() {
		// Create the task instance ids array
		var taskInstanceIds []int

		// Get the task instance ids
		for _, taskInstance := range taskInstances {
			taskInstanceIds = append(taskInstanceIds, taskInstance.ID)
		}

		// Update the task instances
		h.DB.
			Model(&models.TaskInstance{}).
			Where("id IN ?", taskInstanceIds).
			Update("status", models.TaskInstanceAcknowledged).
			Update("acknowledged_at", time.Now())
	}()

	c.JSON(http.StatusOK, taskInstances)
}
