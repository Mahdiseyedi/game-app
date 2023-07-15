package scheduler

import (
	"game-app/param"
	"game-app/service/matchingservice"
	"github.com/go-co-op/gocron"
	_ "github.com/go-co-op/gocron"
	"sync"
	"time"
)

type Scheduler struct {
	sch      *gocron.Scheduler
	matchSvc matchingservice.Service
}

func New(matchSvc matchingservice.Service) Scheduler {
	return Scheduler{
		matchSvc: matchSvc,
		sch:      gocron.NewScheduler(time.UTC)}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	s.sch.Every(3).Second().Do(s.MatchWaitedUsers)

	s.sch.StartAsync()

	<-done
	//here we wait to finish job
	s.sch.Stop()
}

func (s Scheduler) MatchWaitedUsers() {
	s.matchSvc.MatchWaitedUsers(param.MatchWaitedUsersRequest{})
}
