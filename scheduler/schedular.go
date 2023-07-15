package scheduler

import (
	"context"
	"fmt"
	"game-app/param"
	"game-app/service/matchingservice"
	"github.com/go-co-op/gocron"
	_ "github.com/go-co-op/gocron"
	"sync"
	"time"
)

type Config struct {
	MatchWaitedUsersIntervalInSeconds int `koanf:"match_waited_users_interval_in_seconds"`
}

type Scheduler struct {
	sch      *gocron.Scheduler
	matchSvc matchingservice.Service
	config   Config
}

func New(config Config, matchSvc matchingservice.Service) Scheduler {
	return Scheduler{
		config:   config,
		matchSvc: matchSvc,
		sch:      gocron.NewScheduler(time.UTC)}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	//TODO - add lock or distributed lock here
	defer wg.Done()

	s.sch.Every(3).Second().Do(s.MatchWaitedUsers)

	s.sch.Every(s.config.MatchWaitedUsersIntervalInSeconds).Second().Do(s.MatchWaitedUsers)

	s.sch.StartAsync()

	<-done
	//here we wait to finish job
	s.sch.Stop()
}

func (s Scheduler) MatchWaitedUsers() {
	ctx, cancel := context.WithTimeout(context.Background(),
		2*time.Minute)
	defer cancel()

	//get lock(turn)
	_, err := s.matchSvc.MatchWaitedUsers(ctx, param.MatchWaitedUsersRequest{})
	if err != nil {
		//TODO - logging err
		//TODO - update metrics
		fmt.Println("matchSvc.MatchwaitedUsers error: ", err)
	}
	//free lock(turn)
}
