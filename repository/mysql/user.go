package mysql

import (
	"database/sql"
	"fmt"
	"game-app/entity/user"
	"game-app/pkg/errmsg"
	"game-app/pkg/richerror"
)

func (d *MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"
	row := d.db.QueryRow(`select * from users where phone_number=?`, phoneNumber)
	_, err := ScanUser(row)
	if err == sql.ErrNoRows {
		return true, nil
	} else if err != nil {
		return false, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).
			WithKind(richerror.KindUnexpected)
	}

	return false, err
}

func (d *MySQLDB) Register(u user.User) (user.User, error) {
	res, err := d.db.Exec(`insert into users(name, phone_number,password) values(?,?,?)`, u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return user.User{}, fmt.Errorf("cant inseret into MySQLDB, %w", err)
	}

	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d *MySQLDB) GetUserByPhoneNumber(phoneNumber string) (user.User, error) {
	const op = "mysql.GetUserByPhoneNumber"

	row := d.db.QueryRow(`select * from users where phone_number =?`, phoneNumber)
	usr, err := ScanUser(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
		}
		return user.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
	}

	return usr, nil
}

func (d *MySQLDB) GetUserByID(userID uint) (user.User, error) {
	const op = "mysql.GetUserByID"
	row := d.db.QueryRow(`select * from users where id=?`, userID)
	u, err := ScanUser(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
		}
		return user.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)

	}

	return u, nil
}

func ScanUser(row *sql.Row) (user.User, error) {
	var createdAt []uint8
	u := user.User{}
	err := row.Scan(&u.ID, &u.Name, &u.PhoneNumber, &createdAt, &u.Password)
	return u, err
}
