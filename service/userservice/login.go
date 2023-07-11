package userservice

import (
	"fmt"
	"game-app/dto"
	"game-app/pkg/errmsg"
	"game-app/pkg/hash"
	"game-app/pkg/richerror"
)

func (s Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	const op = "Userservice.Login"

	reqUser, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return dto.LoginResponse{},
			richerror.New(op).WithErr(err).WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	if !exist {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	if hash.GetMd5Hash(req.Password) != reqUser.Password {
		return dto.LoginResponse{},
			richerror.New(op).WithErr(err).WithKind(richerror.KindForbidden).
				WithMessage(errmsg.ErrorMsgWrongPassword)
	}

	accessToken, taErr := s.auth.CreateAccessToken(reqUser)
	if taErr != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error: %w", taErr)
	}
	refreshToken, trErr := s.auth.CreateRefreshToken(reqUser)
	if trErr != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error: %w", trErr)
	}

	return dto.LoginResponse{User: dto.UserInfo{
		ID:          reqUser.ID,
		Name:        reqUser.Name,
		PhoneNumber: reqUser.PhoneNumber,
	},
		Tokens: dto.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
