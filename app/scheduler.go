package app

import (
	"time"
)

func ScheduleTask(duration int, callback func()) {
	ticker := time.NewTicker(time.Duration(duration) * time.Minute)

	go func() {
		for {
			<-ticker.C
			callback()
		}
	}()
}
