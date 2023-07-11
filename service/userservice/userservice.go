package userservice

import (
	"game-app/entity/user"
)

type Repository interface {
	Register(u user.User) (user.User, error)
	GetUserByPhoneNumber(phoneNumber string) (user.User, bool, error)
	GetUserByID(userID uint) (user.User, error)
}

type AuthGenerator interface {
	CreateRefreshToken(user user.User) (string, error)
	CreateAccessToken(user user.User) (string, error)
}

type Service struct {
	auth AuthGenerator
	repo Repository
}

func New(repo Repository, authGenerator AuthGenerator) Service {
	return Service{repo: repo, auth: authGenerator}
}
