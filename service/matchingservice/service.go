package matchingservice

import (
	"context"
	"fmt"
	"game-app/entity/category"
	"game-app/entity/player"
	"game-app/entity/waiting_member"
	"game-app/param"
	"game-app/pkg/richerror"
	"game-app/pkg/timestamp"
	"github.com/labstack/gommon/log"
	"github.com/thoas/go-funk"
	_ "github.com/thoas/go-funk"
	"sync"
	"time"
)

type Repo interface {
	AddToWaitingList(userID uint, category category.Category) error
	GetWaitingListByCategory(ctx context.Context,
		category category.Category) ([]waiting_member.WaitingMember,
		error)
}

type PresenceClient interface {
	GetPresence(ctx context.Context, request param.GetPresenceRequest) (param.GetPresenceResponse, error)
}

type Config struct {
	WaitingTimeout time.Duration `koanf:"waiting_timeout"`
}

type Service struct {
	config         Config
	repo           Repo
	presenceClient PresenceClient
}

func New(config Config, repo Repo, presenceClient PresenceClient) Service {
	return Service{config: config, repo: repo, presenceClient: presenceClient}
}

func (s Service) AddToWaitingList(req param.AddToWaitingListRequest) (
	param.AddToWaitingListResponse, error) {
	const op = richerror.Op("matchingservice.AddToWaitingList")

	// add user to the waiting list for the given category if not exist
	err := s.repo.AddToWaitingList(req.UserID, req.Category)
	fmt.Println("matchingservice.AddToWaitingList: ", err)
	if err != nil {
		return param.AddToWaitingListResponse{},
			richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.AddToWaitingListResponse{Timeout: s.config.WaitingTimeout}, nil
}

func (s Service) MatchWaitedUsers(ctx context.Context, _ param.MatchWaitedUsersRequest) (
	param.MatchWaitedUsersResponse, error) {
	const op = richerror.Op("matchingservice.MatchWaitedUsers")

	var wg sync.WaitGroup
	for _, category := range category.CategoryList() {
		wg.Add(1)
		go s.match(ctx, category, &wg)
	}
	wg.Wait()

	return param.MatchWaitedUsersResponse{}, nil
}

func (s Service) match(ctx context.Context, category category.Category, wg *sync.WaitGroup) {
	const op = richerror.Op("matchingservice.match")

	defer wg.Done()

	list, err := s.repo.GetWaitingListByCategory(ctx, category)
	if err != nil {
		//TODO - log Error
		//TODO - update Metrics
		log.Errorf("GetWaitingListByCategory, err: %v\n", err)
		return
	}

	userIDs := make([]uint, 0)
	for _, l := range list {
		userIDs = append(userIDs, l.UserID)
	}

	if len(userIDs) < 2 {
		return
	}

	PresenceList, err := s.presenceClient.GetPresence(ctx, param.GetPresenceRequest{UserIDs: userIDs})
	if err != nil {
		//TODO - log Error
		//TODO - update Metrics

		log.Errorf("GetWaitingListByCategory presenceClient.GetPresence, err: %v\n", err)
		return
	}

	presenceUserIDs := make([]uint, len(list))
	for _, l := range PresenceList.Items {
		presenceUserIDs = append(presenceUserIDs, l.UserID)
	}

	var finalList = make([]waiting_member.WaitingMember, 0)
	for _, l := range list {
		if funk.ContainsUInt(presenceUserIDs, l.UserID) &&
			l.Timestamp < timestamp.Add(-20*time.Second) {
			finalList = append(finalList, l)
		} else {
			//TODO - remove from waiting list
		}
	}

	for i := 0; i < len(finalList)-1; i += 2 {
		mu := player.MatchedUser{
			Category: category,
			UserIDs: []uint{
				finalList[i].UserID,
				finalList[i+1].UserID,
			},
		}

		fmt.Println("mu: ", mu)
		//TODO - publish a new event for mu
		//TODO - remove mu users from waiting list
	}
}
