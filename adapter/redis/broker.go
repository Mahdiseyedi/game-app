package redis

import (
	"context"
	"game-app/entity/event"
	"github.com/labstack/gommon/log"
	"time"
)

func (a Adapter) Publish(event event.Event, payload string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	if err := a.client.Publish(ctx, string(event), payload).Err(); err != nil {
		log.Errorf("publish err: %v\n", err)
		// TODO - log
		// TODO - update metrics
	}

	// TODO - update metrics
}
