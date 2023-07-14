package authorizationservice

import (
	"game-app/entity/permission"
	"game-app/entity/role"
	"game-app/pkg/richerror"
)

type Repository interface {
	GetUserPermissionTitles(userID uint, role role.Role) ([]permission.PermissionTitle, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) CheckAccess(userID uint, role role.Role,
	permissions ...permission.PermissionTitle) (bool, error) {
	const op = "authorizeservice.CheckAccess"

	PermissionTitle, err := s.repo.GetUserPermissionTitles(userID, role)
	if err != nil {
		return false, richerror.New(op).WithErr(err)
	}

	for _, pt := range PermissionTitle {
		for _, p := range permissions {
			if p == pt {
				return true, nil
			}
		}
	}

	return false, nil
}
