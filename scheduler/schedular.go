package scheduler

import (
	"fmt"
	"time"
)

type Scheduler struct{}

func New() Scheduler {
	return Scheduler{}
}

func (s Scheduler) Start(done <-chan bool) {
	fmt.Println("scheduler started")
	for {
		select {
		case d := <-done:
			//here we wait for finish job
			fmt.Println("exiting...", d)
		default:
			now := time.Now()
			fmt.Println("Scheduler now: ", now)
			time.Sleep(3 * time.Second)
		}
	}
}
