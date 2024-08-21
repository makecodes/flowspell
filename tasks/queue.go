package tasks

import (
	"flowspell/models"
	"fmt"

	"gorm.io/gorm"
)

type TaskInstanceHandler struct {
	DB *gorm.DB
}

func Queue() error {
	fmt.Println("Queueing tasks")
	return nil
}

func QueueCleanup() error {
	h := TaskInstanceHandler{}
	tasks, err := models.GetAcknowledgedTasks(h.DB)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		fmt.Println("Cleaning up task", task.ID)
	}
	return nil
}
