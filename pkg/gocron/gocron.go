package gocron

import (
	"time"

	"github.com/go-co-op/gocron"
)

func NewSchedule() *gocron.Scheduler {
	return gocron.NewScheduler(time.Local)
}
