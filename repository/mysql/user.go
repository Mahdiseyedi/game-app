package mysql

import (
	"database/sql"
	"fmt"
	"game-app/entity/user"
)

func (d *DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	u := user.User{}
	var createdAt []uint8

	row := d.db.QueryRow(`select * from users where phone_number=?`, phoneNumber)
	err := row.Scan(&u.ID, &u.Name, &u.PhoneNumber, &createdAt)
	if err == sql.ErrNoRows {
		return true, err
	} else if err != nil {
		return false, fmt.Errorf("some thing went wrong, %w", err)
	}

	return false, err
}

func (d *DB) RegisterUser(u user.User) (user.User, error) {
	res, err := d.db.Exec(`insert into users(name, phone_number) values(?,?)`, u.Name, u.PhoneNumber)
	if err != nil {
		return user.User{}, fmt.Errorf("cant inseret into DB, %w", err)
	}

	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}
