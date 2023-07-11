package userservice

import (
	"game-app/dto"
	"game-app/pkg/richerror"
)

func (s Service) Profile(req dto.ProfileRequest) (dto.ProfileResponse, error) {
	const op = "userService.Profile"

	userProfile, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return dto.ProfileResponse{},
			richerror.New(op).WithErr(err).WithMeta(map[string]interface{}{"req": req})
	}

	return dto.ProfileResponse{Name: userProfile.Name}, nil
}
