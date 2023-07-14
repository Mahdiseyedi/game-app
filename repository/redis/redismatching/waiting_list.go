package redismatching

import (
	"context"
	"fmt"
	"game-app/entity/category"
	"game-app/pkg/richerror"
	"github.com/redis/go-redis/v9"
	"time"
)

//TODO - add to config in use-case layer
const WaitingListPrefix = "waitinglist"

func (d DB) AddToWaitingList(userID uint, category category.Category) error {
	const op = richerror.Op("redismatching.AddToWaitingList")

	t, err := d.adapter.Client().
		ZAdd(context.Background(),
			fmt.Sprintf("%s:%s", WaitingListPrefix, category),
			redis.Z{Score: float64(time.Now().UnixMicro()),
				Member: fmt.Sprintf("%d", userID),
			}).Result()

	fmt.Println("redismatching.AddToWaitingList t: ", t)
	fmt.Println("redismatching.AddToWaitingList err: ", err)
	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}
