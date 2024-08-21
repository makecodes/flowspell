package cron

import (
	"flowspell/tasks"

	"github.com/robfig/cron/v3"
)

func Start() {
	c := cron.New()
	// c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	// c.AddFunc("@hourly",      func() { fmt.Println("Every hour") })
	c.AddFunc("@every 5s", func() { tasks.Queue() })
	c.Start()
	select {}
}
