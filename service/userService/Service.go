package userService

import (
	"fmt"
	"game-app/entity"
	"game-app/pkg/phoneNumber"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	RegisterUser(u entity.User) (entity.User, error)
}

type Service struct {
	repo Repository
}

type RegisterRequest struct {
	Name        string
	PhoneNumber string
}

type RegisterResponse struct {
	user entity.User
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	//TODO - implementing otp verification for phoneNumber
	if !phoneNumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("this number is not valid")
	}

	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); !isUnique || err != nil {
		if err != nil {
			return RegisterResponse{}, err
		}

		if !isUnique {
			return RegisterResponse{}, fmt.Errorf("phone number is not unique")
		}
	}

	if len(req.Name) <= 3 {
		return RegisterResponse{}, fmt.Errorf("name lenght should grater than 3")
	}
	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	}

	createdUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexcepted Register Error %w", err)
	}

	return RegisterResponse{
		user: createdUser,
	}, nil
}
