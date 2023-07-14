package game

import (
	"game-app/entity/category"
	"game-app/entity/question"
	"time"
)

type Game struct {
	ID           uint
	Category     category.Category
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
	Choice     question.PossibleAnswerChoice
}

func Data() {

}
