package cron

import (
	"flowspell/db"
	"flowspell/tasks"
	"log"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Start() {
	dbConnection, err := db.GetDBConnection()
	if err != nil {
		log.Fatalf("Error connecting to database")
	}

	cronTasksHandler := tasks.TaskInstanceHandler{DB: dbConnection}

	c := cron.New()
	// c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	// c.AddFunc("@hourly",      func() { fmt.Println("Every hour") })
	_, err = c.AddFunc("@every 5s", func() { cronTasksHandler.QueueCleanup() })
	if err != nil {
		return
	}
	c.Start()
	select {}
}
