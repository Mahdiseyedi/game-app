package entity

import "time"

type Game struct {
	ID           uint
	CategoryID   uint
	QuestionsIDs []uint
	PlayerIDs    []uint
	startTime    time.Time
}

type Player struct {
	ID      uint
	UserID  uint
	GameID  uint
	Score   uint
	Answers []PlayerAnswer
}

type PlayerAnswer struct {
	ID         uint
	PlayerID   uint
	QuestionID uint
	Choice     PossibleAnswerChoice
}
