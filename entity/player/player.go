package player

import "game-app/entity/category"

type MatchedUser struct {
	Category category.Category
	UserIDs  []uint
}
