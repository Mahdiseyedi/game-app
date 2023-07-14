package authorizationservice

import (
	"fmt"
	"game-app/entity/permission"
	"game-app/entity/role"
	"game-app/pkg/richerror"
	"reflect"
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
	fmt.Println("authorizationservice.PermissionTitle: ", PermissionTitle)
	fmt.Println("authorizationservice.erroe	:", err)
	if err != nil {
		return false, richerror.New(op).WithErr(err)
	}

	for _, pt := range PermissionTitle {
		fmt.Println("CheckAccess.type pt: ", reflect.TypeOf(pt), "--pt: ", pt)
		for _, p := range permissions {
			fmt.Println("CheckAccess.type p: ", reflect.TypeOf(p), "--p: ", p)
			if p == pt {
				return true, nil
			}
		}
	}

	return false, nil
}
