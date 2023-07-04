package mysql

import (
	"database/sql"
	"fmt"
	"game-app/entity"
)

func (d *DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	u := entity.User{}
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

//select * from users where phone_number=?

func (d *DB) RegisterUser(u entity.User) (entity.User, error) {
	res, err := d.db.Exec(`insert into users(name, phone_number) values(?,?)`, u.Name, u.PhoneNumber)
	if err != nil {
		return entity.User{}, fmt.Errorf("cant inseret into DB, %w", err)
	}

	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}
