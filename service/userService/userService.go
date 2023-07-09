package userService

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

type UserInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type RegisterResponse struct {
	UsrInfo UserInfo `json:"user_info"`
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
	User   UserInfo `json:"user_info"`
	Tokens Tokens   `json:"tokens"`
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

func (s Service) Register(req dto.RegisterRequest) (RegisterResponse, error) {
	//TODO - implementing otp verification for phoneNumber

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
		UsrInfo: UserInfo{
			ID:          createdUser.ID,
			Name:        createdUser.Name,
			PhoneNumber: createdUser.PhoneNumber,
		},
	}, nil
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	const op = "UserService.Login"

	reqUser, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{},
			richerror.New(op).WithErr(err).WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	if !exist {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	if hash.GetMd5Hash(req.Password) != reqUser.Password {
		return LoginResponse{}, fmt.Errorf("...Service: Login failed!...")
	}

	accessToken, taErr := s.auth.CreateAccessToken(reqUser)
	if taErr != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", taErr)
	}
	refreshToken, trErr := s.auth.CreateRefreshToken(reqUser)
	if trErr != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", trErr)
	}

	return LoginResponse{User: UserInfo{
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

func (s Service) GetUserProfile(req UserProfileRequest) (UserProfileResponse, error) {
	const op = "userService.GetUserProfile"

	userProfile, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return UserProfileResponse{},
			richerror.New(op).WithErr(err).WithMeta(map[string]interface{}{"req": req})
	}

	return UserProfileResponse{Name: userProfile.Name}, nil
}

func (s Service) PasswordServiceValidator(req dto.RegisterRequest) (bool, error) {
	//TODO - check the password with regex
	if len(req.Password) < 8 {
		return false, fmt.Errorf("...Validator: Password len most grater than 8...")
	}

	if req.Password == "password" {
		return false, fmt.Errorf("...Validator: so simple...")
	}

	return true, nil
}
