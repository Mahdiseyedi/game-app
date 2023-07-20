package broker

import "game-app/entity/event"

type Publisher interface {
	Publish(event event.Event, payload string)
}
