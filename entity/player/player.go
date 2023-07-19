package player

import "game-app/entity/category"

type MatchedUsers struct {
	Category category.Category
	UserIDs  []uint
}
