package matchingservice

import (
	"context"
	"fmt"
	"game-app/contract/broker"
	"game-app/entity/category"
	"game-app/entity/event"
	"game-app/entity/player"
	"game-app/entity/waiting_member"
	"game-app/param"
	"game-app/pkg/protobufencoder"
	"game-app/pkg/richerror"
	"game-app/pkg/timestamp"
	"github.com/labstack/gommon/log"
	_ "github.com/thoas/go-funk"
	"sync"
	"time"
)

//TODO - add ctx to all repo and use-case method if need
type Repo interface {
	AddToWaitingList(userID uint, category category.Category) error
	GetWaitingListByCategory(ctx context.Context,
		category category.Category) ([]waiting_member.WaitingMember,
		error)
	RemoveUsersFromWaitingList(category category.Category, userIDs []uint)
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
	pub            broker.Publisher
}

func New(config Config, repo Repo, presenceClient PresenceClient, pub broker.Publisher) Service {
	return Service{config: config, repo: repo, presenceClient: presenceClient, pub: pub}
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

	PresenceList, err := s.presenceClient.GetPresence(ctx,
		param.GetPresenceRequest{UserIDs: userIDs})
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

	var toBeRemovedUsers = make([]uint, 0)
	var finalList = make([]waiting_member.WaitingMember, 0)
	for _, l := range list {
		lastOnLineTimestamp, ok := getPresenceItem(PresenceList, l.UserID)
		if ok && lastOnLineTimestamp > timestamp.Add(-60*time.Second) &&
			l.Timestamp > timestamp.Add(-300*time.Second) {
			finalList = append(finalList, l)
		} else {
			//remove from waiting list
			toBeRemovedUsers = append(toBeRemovedUsers, l.UserID)
		}
	}

	go s.repo.RemoveUsersFromWaitingList(category, toBeRemovedUsers)

	matchedUsersToBeRemoved := make([]uint, 0)
	for i := 0; i < len(finalList)-1; i += 2 {
		mu := player.MatchedUsers{
			Category: category,
			UserIDs: []uint{
				finalList[i].UserID,
				finalList[i+1].UserID,
			},
		}
		fmt.Println("mu: ", mu)

		go s.pub.Publish(event.MatchingUsersMatchedEvent,
			protobufencoder.EncodeMatchingUsersMatchedEvent(mu))

		//remove mu users from waiting list
		matchedUsersToBeRemoved = append(matchedUsersToBeRemoved, mu.UserIDs...)
	}

	go s.repo.RemoveUsersFromWaitingList(category, matchedUsersToBeRemoved)
}

func getPresenceItem(presenceList param.GetPresenceResponse, userID uint) (int64, bool) {
	for _, item := range presenceList.Items {
		if item.UserID == userID {
			return item.Timestamp, true
		}
	}

	return 0, false
}
