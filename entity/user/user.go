package user

import "game-app/entity/role"

type User struct {
	ID          uint
	Name        string
	PhoneNumber string
	Password    string
	Role        role.Role
}
