package waiting_member

import "game-app/entity/category"

type WaitingMember struct {
	UserID    uint
	Timestamp int64
	Category  category.Category
}
