package scheduler

import (
	"log"
	"time"
	"video-conf/database"
	"video-conf/database/models"

	"github.com/go-co-op/gocron/v2"
)

func eventsStatusScheduler(scheduler gocron.Scheduler) {
	job, err := scheduler.NewJob(
		gocron.DurationJob(
			time.Minute,
		),
		gocron.NewTask(
			func() {
				var events []models.Event

				DB := database.Connection()

				DB.Table("events").
					Where("status != ?", models.EVENT_FINISHED_STATUS).
					Where("start_at <= ?", time.Now()).
					Find(&events)

				for _, event := range events {
					endTime := event.StartAt.Add(time.Minute * time.Duration(event.Duration))

					if time.Now().After(endTime) {
						event.Status = models.EVENT_FINISHED_STATUS
						DB.Save(&event)
						log.Println("Event status value changed:", event.Title)
					}
				}

				if len(events) == 0 {
					log.Println("Events status scheduler did not find any event")
				}
			},
		),
	)

	if err != nil {
		log.Fatalf("Error running events status job")
	}

	log.Println("Events status job is running:", job.ID())
}

func InitScheduler() {
	scheduler, err := gocron.NewScheduler()

	if err != nil {
		log.Fatalf("Error running scheduler")
	}

	eventsStatusScheduler(scheduler)

	scheduler.Start()
}
