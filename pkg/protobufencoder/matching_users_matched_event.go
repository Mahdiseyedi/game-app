package protobufencoder

import (
	"encoding/base64"
	"game-app/contract/golang/matching"
	"game-app/entity/category"
	"game-app/entity/player"
	"game-app/pkg/slice"
	"github.com/golang/protobuf/proto"
)

func EncodeMatchingUsersMatchedEvent(mu player.MatchedUsers) string {
	pbMu := matching.MatchedUsers{
		UserIDs:  slice.MapFromUintToUint64(mu.UserIDs),
		Category: string(mu.Category),
	}

	payload, err := proto.Marshal(&pbMu)
	if err != nil {
		// TODO - log error
		// TODO - update metrics
		return ""
	}

	return base64.StdEncoding.EncodeToString(payload)
}

func DecodeMatchingUsersMatchedEvent(data string) player.MatchedUsers {
	payload, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		// TODO - log
		// TODO - update metrics
		return player.MatchedUsers{}
	}

	pbMu := matching.MatchedUsers{}
	if err := proto.Unmarshal(payload, &pbMu); err != nil {
		// TODO - log
		// TODO - update metrics
		return player.MatchedUsers{}
	}

	return player.MatchedUsers{
		Category: category.Category(pbMu.Category),
		UserIDs:  slice.MapFromUint64ToUint(pbMu.UserIDs),
	}
}
