package userservice

import (
	"context"
	"game-app/param"
	"game-app/pkg/richerror"
)

func (s Service) Profile(ctx context.Context, req param.ProfileRequest) (param.ProfileResponse, error) {
	const op = "userService.Profile"

	userProfile, err := s.repo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return param.ProfileResponse{},
			richerror.New(op).WithErr(err).WithMeta(map[string]interface{}{"req": req})
	}

	return param.ProfileResponse{Name: userProfile.Name}, nil
}
