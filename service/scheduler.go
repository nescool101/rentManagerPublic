package service

import (
	"log"

	"github.com/robfig/cron/v3"
)

// StartScheduler initializes and starts the cron scheduler.
func StartScheduler() {
	c := cron.New()
	_, err := c.AddFunc("@monthly", func() {
		log.Println("Running monthly rent reminders...")
		//NotifyAll()
	})
	if err != nil {
		log.Fatalf("Error adding scheduler function: %v", err)
	}
	c.Start()

	// Prevent the scheduler from exiting
	select {}
}
