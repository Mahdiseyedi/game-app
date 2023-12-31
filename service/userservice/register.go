package userservice

import (
	"game-app/entity/role"
	users "game-app/entity/user"
	"game-app/param"
	"game-app/pkg/errmsg"
	"game-app/pkg/hash"
	"game-app/pkg/richerror"
)

func (s Service) Register(req param.RegisterRequest) (param.RegisterResponse, error) {
	const op = "userservice.Register"
	//TODO - implementing otp verification for phoneNumber

	//TODO - replace md5 with bcrypt
	user := users.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    hash.GetMd5Hash(req.Password),
		Role:        role.UserRole,
	}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return param.RegisterResponse{}, richerror.New(op).WithErr(err).
			WithKind(richerror.KindUnexpected).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong)
	}

	return param.RegisterResponse{User: param.UserInfo{
		ID:          createdUser.ID,
		Name:        createdUser.Name,
		PhoneNumber: createdUser.PhoneNumber,
	}}, nil
}
