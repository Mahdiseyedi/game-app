package userService

import (
	"fmt"
	"game-app/entity/user"
	"game-app/pkg/hash"
	"game-app/pkg/phoneNumber"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	RegisterUser(u user.User) (user.User, error)
	GetUserByPhoneNumber(phoneNumber string) (user.User, error)
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

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProfileRequest struct {
	UserID uint `json:"user_id"`
}

type UserProfileResponse struct {
	Name string
}

func New(repo Repository, auth AuthGenerator) Service {
	return Service{repo: repo, auth: auth}
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	//TODO - implementing otp verification for phoneNumber
	if res, err := s.PhoneNumberServiceValidator(req); !res {
		return RegisterResponse{}, err
	}

	if res, err := s.NameServiceValidator(req); !res {
		return RegisterResponse{}, err
	}

	if res, err := s.PasswordServiceValidator(req); !res {
		return RegisterResponse{}, err
	}

	//TODO - replace md5 with bcrypt
	hashedPassword := hash.GetMd5Hash(req.Password)

	user := user.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,
	}

	createdUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("...Repository: Register repository Error %w", err)
	}

	return RegisterResponse{
		ID:          createdUser.ID,
		Name:        createdUser.Name,
		PhoneNumber: createdUser.PhoneNumber,
	}, nil
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	reqUser, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, err
	}

	if hash.GetMd5Hash(req.Password) != reqUser.Password {
		return LoginResponse{}, fmt.Errorf("...Service: Login failed!...")
	}

	accessToken, tErr := s.auth.CreateAccessToken(reqUser)
	if tErr != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", tErr)
	}
	refreshToken, tErr := s.auth.CreateRefreshToken(reqUser)
	if tErr != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", tErr)
	}

	return LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s Service) GetUserProfile(req UserProfileRequest) (UserProfileResponse, error) {
	userProfile, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return UserProfileResponse{}, err
	}

	return UserProfileResponse{Name: userProfile.Name}, nil
}

func (s Service) PhoneNumberServiceValidator(req RegisterRequest) (bool, error) {
	if !phoneNumber.IsValid(req.PhoneNumber) {
		return false, fmt.Errorf("...Validator: this number is not valid...")
	}

	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); !isUnique {
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

func (s Service) PasswordServiceValidator(req RegisterRequest) (bool, error) {
	//TODO - check the password with regex
	if len(req.Password) < 8 {
		return false, fmt.Errorf("...Validator: Password len most grater than 8...")
	}

	if req.Password == "password" {
		return false, fmt.Errorf("...Validator: so simple...")
	}

	return true, nil
}
