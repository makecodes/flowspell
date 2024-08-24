package tasks

import (
	"flowspell/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TaskInstanceHandler struct {
	DB *gorm.DB
}

func EnqueueTask() error {
	fmt.Println("EnqueueTask")
	return nil
}

func (h *TaskInstanceHandler) QueueCleanup() {
	var tasks []models.TaskInstance
	err := h.DB.
		Where("status = ? AND acknowledged_at < ?", models.TaskInstanceAcknowledged, time.Now().
			Add(-1*time.Minute)).
		Find(&tasks).
		Error

	if err != nil {
		return
	}

	for _, task := range tasks {
		// Each task with status "acknowledged" will be cleaned up status TaskInstanceStatusNotStarted
		task.Status = models.TaskInstanceStatusNotStarted
		if err := h.DB.Save(&task).Error; err != nil {
			return
		}

		fmt.Println("Cleaning up task", task.ID, task.Name)
	}
}
