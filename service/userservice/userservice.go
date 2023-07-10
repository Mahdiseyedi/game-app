package userservice

import (
	"fmt"
	"game-app/dto"
	"game-app/entity/user"
	"game-app/pkg/hash"
	"game-app/pkg/richerror"
)

type Repository interface {
	RegisterUser(u user.User) (user.User, error)
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

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	User   dto.UserInfo `json:"user_info"`
	Tokens Tokens       `json:"tokens"`
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

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	//TODO - implementing otp verification for phoneNumber

	//TODO - replace md5 with bcrypt
	user := user.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    hash.GetMd5Hash(req.Password),
	}

	createdUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("...Repository: Register repository Error %w", err)
	}

	return dto.RegisterResponse{User: dto.UserInfo{
		ID:          createdUser.ID,
		Name:        createdUser.Name,
		PhoneNumber: createdUser.PhoneNumber,
	}}, nil
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	const op = "Userservice.userLogin"

	reqUser, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{},
			richerror.New(op).WithErr(err).WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	if !exist {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	if hash.GetMd5Hash(req.Password) != reqUser.Password {
		return LoginResponse{}, fmt.Errorf("...Service: userLogin failed!...")
	}

	accessToken, taErr := s.auth.CreateAccessToken(reqUser)
	if taErr != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", taErr)
	}
	refreshToken, trErr := s.auth.CreateRefreshToken(reqUser)
	if trErr != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", trErr)
	}

	return LoginResponse{User: dto.UserInfo{
		ID:          reqUser.ID,
		Name:        reqUser.Name,
		PhoneNumber: reqUser.PhoneNumber,
	},
		Tokens: Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s Service) Profile(req UserProfileRequest) (UserProfileResponse, error) {
	const op = "userService.Profile"

	userProfile, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return UserProfileResponse{},
			richerror.New(op).WithErr(err).WithMeta(map[string]interface{}{"req": req})
	}

	return UserProfileResponse{Name: userProfile.Name}, nil
}

func (s Service) PasswordServiceValIDator(req dto.RegisterRequest) (bool, error) {
	//TODO - check the password with regex
	if len(req.Password) < 8 {
		return false, fmt.Errorf("...ValIDator: Password len most grater than 8...")
	}

	if req.Password == "password" {
		return false, fmt.Errorf("...ValIDator: so simple...")
	}

	return true, nil
}
