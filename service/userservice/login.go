package userservice

import (
	"fmt"
	"game-app/param"
	"game-app/pkg/errmsg"
	"game-app/pkg/hash"
	"game-app/pkg/richerror"
)

func (s Service) Login(req param.LoginRequest) (param.LoginResponse, error) {
	const op = "userService.Login"

	reqUser, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return param.LoginResponse{},
			richerror.New(op).WithErr(err).WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	if hash.GetMd5Hash(req.Password) != reqUser.Password {
		return param.LoginResponse{},
			richerror.New(op).WithErr(err).WithKind(richerror.KindForbidden).
				WithMessage(errmsg.ErrorMsgWrongPassword)
	}

	accessToken, taErr := s.auth.CreateAccessToken(reqUser)
	if taErr != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error: %w", taErr)
	}
	refreshToken, trErr := s.auth.CreateRefreshToken(reqUser)
	if trErr != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error: %w", trErr)
	}

	return param.LoginResponse{User: param.UserInfo{
		ID:          reqUser.ID,
		Name:        reqUser.Name,
		PhoneNumber: reqUser.PhoneNumber,
	},
		Tokens: param.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
