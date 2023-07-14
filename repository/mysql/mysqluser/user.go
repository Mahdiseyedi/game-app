package mysqluser

import (
	"database/sql"
	"game-app/entity/role"
	"game-app/entity/user"
	"game-app/pkg/errmsg"
	"game-app/pkg/richerror"
	"game-app/repository/mysql"
)

func (d *DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"
	row := d.conn.Conn().QueryRow(`select * from users where phone_number = ?`, phoneNumber)

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

func (d *DB) Register(u user.User) (user.User, error) {
	const op = "mysql.Register"

	res, err := d.conn.Conn().Exec(`insert into users(name, phone_number, password, role) values(?,?,?,?)`,
		u.Name, u.PhoneNumber, u.Password, u.Role.String())
	if err != nil {
		return user.User{}, richerror.New(op).WithErr(err).
			WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrorMsgCantInsertUserIntoDatabase)
	}

	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d *DB) GetUserByPhoneNumber(phoneNumber string) (user.User, error) {
	const op = "mysql.GetUserByPhoneNumber"

	row := d.conn.Conn().QueryRow(`select * from users where phone_number = ?`, phoneNumber)
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

func (d *DB) GetUserByID(userID uint) (user.User, error) {
	const op = "mysql.GetUserByID"
	row := d.conn.Conn().QueryRow(`select * from users where id = ?`, userID)
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

func ScanUser(scanner mysql.Scanner) (user.User, error) {
	var createdAt []uint8
	var user user.User
	var roleStr string

	err := scanner.Scan(&user.ID, &user.Name, &user.PhoneNumber,
		&createdAt, &user.Password, &roleStr)

	user.Role = role.MapToRoleEntity(roleStr)

	return user, err
}
