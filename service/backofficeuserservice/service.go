package backofficeuserservice

import (
	"fmt"
	"game-app/entity/role"
	"game-app/entity/user"
)

type Service struct{}

func New() Service {
	return Service{}
}

func (s Service) ListAllUsers() ([]user.User, error) {
	//TODO - implement me
	list := make([]user.User, 0)

	list = append(list, user.User{
		ID:          0,
		Name:        "fake",
		PhoneNumber: "fake",
		Password:    "fake",
		Role:        role.AdminRole,
	})
	fmt.Println("backofficeuserservice.ListAllUsers: ", list)

	return list, nil
}
