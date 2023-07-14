package param

import (
	"game-app/entity/category"
	"time"
)

type AddToWaitingListRequest struct {
	UserID   uint              `json:"user_id"`
	Category category.Category `json:"category"`
}

type AddToWaitingListResponse struct {
	Timeout time.Duration `json:"timeout_in_nanoseconds"`
}
