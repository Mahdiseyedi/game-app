package main

import (
	"context"
	"fmt"
	"game-app/adapter/redis"
	"game-app/config"
	"game-app/entity/category"
	"game-app/entity/event"
	"game-app/entity/player"
	"game-app/pkg/protobufencoder"
)

func main() {
	cfg := config.Load("../../../config.yml")

	redisAdapter := redis.New(cfg.Redis)
	topic := event.MatchingUsersMatchedEvent

	mu := player.MatchedUsers{
		Category: category.FootballCategory,
		UserIDs:  []uint{1, 4},
	}

	payload := protobufencoder.EncodeMatchingUsersMatchedEvent(mu)
	if err := redisAdapter.Client().Publish(context.Background(), string(topic), payload).Err(); err != nil {
		panic(fmt.Sprintf("publish err: %v", err))
	}
}
