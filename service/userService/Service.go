package userService

import (
	"fmt"
	"game-app/entity/user"
	"game-app/pkg/phoneNumber"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	RegisterUser(u user.User) (user.User, error)
}

type Service struct {
	repo Repository
}

type RegisterRequest struct {
	Name        string
	PhoneNumber string
}

type RegisterResponse struct {
	user user.User
}

func New(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	//TODO - implementing otp verification for phoneNumber
	if res, err := s.PhoneNumberServiceValidator(req); !res {
		return RegisterResponse{}, err
	}

	if res, err := s.NameServiceValidator(req); !res {
		return RegisterResponse{}, err
	}

	user := user.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	}

	createdUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("...Repository: Register repository Error %w", err)
	}

	return RegisterResponse{
		user: createdUser,
	}, nil
}

func (s Service) PhoneNumberServiceValidator(req RegisterRequest) (bool, error) {
	if !phoneNumber.IsValid(req.PhoneNumber) {
		return false, fmt.Errorf("...Validator: this number is not valid...")
	}

	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); !isUnique || err == nil {
		if err != nil {
			return false, err
		}

		if !isUnique {
			return false, fmt.Errorf("...Validator: phone number is not unique...")
		}
	}
	return true, nil
}

func (s Service) NameServiceValidator(req RegisterRequest) (bool, error) {
	if len(req.Name) <= 3 {
		return false, fmt.Errorf("...Validator: name lenght should grater than 3")
	}

	//this if statement just for keep validator structure regular
	if req.Name == "userName" {
		return false, fmt.Errorf("...Validator: username cant be \"userName\"")
	}

	return true, nil
}
