package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"game-app/adapter/redis"
	"game-app/config"
	"game-app/contract/golang/matching"
	"game-app/entity/category"
	"game-app/entity/player"
	"game-app/pkg/slice"
	"github.com/golang/protobuf/proto"
)

func main() {
	cfg := config.Load("../../../config.yml")

	redisAdapter := redis.New(cfg.Redis)
	topic := "matching.users_matched"

	mu := player.MatchedUser{Category: category.FootballCategory,
		UserIDs: []uint{1, 4},
	}

	pbMu := matching.MatchedUsers{
		Category: string(mu.Category),
		UserIDs:  slice.MapFromUintToUint64(mu.UserIDs),
	}
	payload, err := proto.Marshal(&pbMu)
	if err != nil {
		panic(fmt.Sprintf("payload err: %v", err))
	}

	payloadStr := base64.StdEncoding.EncodeToString(payload)

	if err := redisAdapter.Client().Publish(context.Background(), topic, payloadStr).Err(); err != nil {
		panic(fmt.Sprintf("publish err: %v", err))
	}
}
