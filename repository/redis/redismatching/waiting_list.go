package redismatching

import (
	"context"
	"fmt"
	"game-app/entity/category"
	"game-app/entity/waiting_member"
	"game-app/pkg/richerror"
	"game-app/pkg/timestamp"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

//TODO - add to config in use-case layer
const WaitingListPrefix = "waitinglist"

func (d DB) AddToWaitingList(userID uint, category category.Category) error {
	const op = richerror.Op("redismatching.AddToWaitingList")

	t, err := d.adapter.Client().
		ZAdd(context.Background(),
			fmt.Sprintf("%s:%s", WaitingListPrefix, category),
			redis.Z{Score: float64(timestamp.Now()),
				Member: fmt.Sprintf("%d", userID),
			}).Result()

	fmt.Println("redismatching.AddToWaitingList t: ", t)
	fmt.Println("redismatching.AddToWaitingList err: ", err)
	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (d DB) GetWaitingListByCategory(ctx context.Context, category category.Category) (
	[]waiting_member.WaitingMember, error) {
	const op = richerror.Op("redismatching.GetWaitingListByCategory")
	min := fmt.Sprintf("%d", timestamp.Add(-2*time.Hour))
	max := strconv.Itoa(int(timestamp.Now()))

	list, err := d.adapter.Client().ZRevRangeByScoreWithScores(ctx,
		getCategoryKey(category), &redis.ZRangeBy{
			Min:    min,
			Max:    max,
			Offset: 0,
			Count:  0,
		}).Result()

	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithKind(richerror.KindUnexpected)
	}

	var result = make([]waiting_member.WaitingMember, 0)

	for _, l := range list {
		userID, _ := strconv.Atoi(l.Member.(string))

		result = append(result, waiting_member.WaitingMember{
			UserID:    uint(userID),
			Timestamp: int64(l.Score),
			Category:  category,
		})
	}

	return result, nil
}

func getCategoryKey(category category.Category) string {
	return fmt.Sprintf("%s:%s", WaitingListPrefix, category)
}
