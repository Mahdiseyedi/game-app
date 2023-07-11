package userservice

import (
	"fmt"
	"game-app/dto"
	users "game-app/entity/user"
	"game-app/pkg/hash"
)

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
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
		return dto.RegisterResponse{}, fmt.Errorf("...Repository: Register repository Error %w", err)
	}

	return dto.RegisterResponse{User: dto.UserInfo{
		ID:          createdUser.ID,
		Name:        createdUser.Name,
		PhoneNumber: createdUser.PhoneNumber,
	}}, nil
}
