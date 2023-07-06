package mysql

import (
	"database/sql"
	"fmt"
	"game-app/entity/user"
)

func (d *DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	row := d.db.QueryRow(`select * from users where phone_number=?`, phoneNumber)
	_, err := ScanUser(row)
	if err == sql.ErrNoRows {
		return true, err
	} else if err != nil {
		return false, fmt.Errorf("...some thing went wrong, %w", err)
	}

	return false, err
}

func (d *DB) RegisterUser(u user.User) (user.User, error) {
	res, err := d.db.Exec(`insert into users(name, phone_number,password) values(?,?,?)`, u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return user.User{}, fmt.Errorf("cant inseret into DB, %w", err)
	}

	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d *DB) GetUserByPhoneNumber(phoneNumber string) (user.User, error) {

	row := d.db.QueryRow(`select * from users where phone_number =?`, phoneNumber)
	u, err := ScanUser(row)

	if err != nil {
		return user.User{}, fmt.Errorf("...no user find with that phone number, %w", err)
	}

	return u, nil
}

func (d *DB) GetUserByID(userID uint) (user.User, error) {
	row := d.db.QueryRow(`select * from users where id=?`, userID)
	u, err := ScanUser(row)

	if err != nil {
		return user.User{}, fmt.Errorf("...no user find with that ID, %w", err)
	}
	return u, nil
}

func ScanUser(row *sql.Row) (user.User, error) {
	var createdAt []uint8
	u := user.User{}
	err := row.Scan(&u.ID, &u.Name, &u.PhoneNumber, &createdAt, &u.Password)

	return u, err
}
