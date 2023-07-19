package event

type Event string

const (
	MatchingUsersMatchedEvent Event = "matching.users_matched"
	ServiceFailedToMatching   Event = "matching.failed_to_matched"
)
