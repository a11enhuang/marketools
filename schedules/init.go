package schedules

import (
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

var instance *Scheduler

func init() {
	instance = &Scheduler{
		scheduler: gocron.NewScheduler(time.Local),
		taskQueue: make(chan func(), 256),
	}
}

func AsyncStart(ctx context.Context) {
	go instance.start(ctx)
}

type Scheduler struct {
	scheduler *gocron.Scheduler
	taskQueue chan func()
}

func (this *Scheduler) start(ctx context.Context) {
	this.scheduler.StartAsync()
	for {
		select {
		case <-ctx.Done():
			log.Println("[Scheduler]任务调度器已结束")
			return
		case task := <-this.taskQueue:
			task()
		}
	}
}

func AddWeekdaysTask(time string, job func()) {
	instance.scheduler.Every(1).
		Monday().Tuesday().Wednesday().Thursday().Friday().
		At(time).
		Do(func() {
			instance.taskQueue <- job
		})
}
