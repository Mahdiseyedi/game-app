package main

import (
	"context"
	"fmt"
	"game-app/adapter/redis"
	"game-app/config"
	"game-app/entity/event"
	"game-app/pkg/protobufencoder"
)

func main() {
	cfg := config.Load("../../../config.yml")

	redisAdapter := redis.New(cfg.Redis)

	topic := event.MatchingUsersMatchedEvent

	subscriber := redisAdapter.Client().Subscribe(context.Background(), string(topic))

	for {
		msg, err := subscriber.ReceiveMessage(context.Background())
		if err != nil {
			panic(fmt.Sprintf("main.subscriber.ReceiveMessage: %v", err))
		}

		switch event.Event(msg.Channel) {
		case topic:
			processUsersMatchedEvent(msg.Channel, msg.Payload)
		default:
			fmt.Println("invalid topic... ", msg.Channel)
		}
	}
}

func processUsersMatchedEvent(topic string, data string) {

	mu := protobufencoder.DecodeMatchingUsersMatchedEvent(data)

	fmt.Println("Received message from " + topic + " topic.")
	fmt.Printf("matched users %+v\n", mu)
}
