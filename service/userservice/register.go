package userservice

import (
	"fmt"
	users "game-app/entity/user"
	"game-app/param"
	"game-app/pkg/hash"
)

func (s Service) Register(req param.RegisterRequest) (param.RegisterResponse, error) {
	//TODO - implementing otp verification for phoneNumber

	//TODO - replace md5 with bcrypt
	user := users.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    hash.GetMd5Hash(req.Password),
	}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return param.RegisterResponse{}, fmt.Errorf("...Repository: Register repository Error %w", err)
	}

	return param.RegisterResponse{User: param.UserInfo{
		ID:          createdUser.ID,
		Name:        createdUser.Name,
		PhoneNumber: createdUser.PhoneNumber,
	}}, nil
}
