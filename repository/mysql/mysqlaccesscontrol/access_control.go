package mysqlaccesscontrol

import (
	"game-app/entity/accesscontrol"
	"game-app/entity/permission"
	"game-app/entity/role"
	"game-app/pkg/errmsg"
	"game-app/pkg/richerror"
	"game-app/pkg/slice"
	"game-app/repository/mysql"
	"strings"
	"time"
)

func (d *DB) GetUserPermissionTitles(userID uint, role role.Role) ([]permission.PermissionTitle, error) {
	const op = "mysql.GetUserPermissionTitles"

	roleACL := make([]accesscontrol.AccessControl, 0)
	rows, err := d.conn.Conn().Query(`select * from access_controls where actor_type = ? and actor_id = ?`,
		accesscontrol.RoleActorType, role)

	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		acl, err := scanAccessControl(rows)

		if err != nil {
			return nil, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
		}

		roleACL = append(roleACL, acl)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	userACL := make([]accesscontrol.AccessControl, 0)

	userRows, err := d.conn.Conn().Query(`select * from access_controls where actor_type = ? actor_id = ?`,
		accesscontrol.UserActorType, userID)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	defer userRows.Close()

	for userRows.Next() {
		acl, err := scanAccessControl(userRows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
		}

		userACL = append(userACL, acl)

	}

	if err := userRows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	//merge ACLs by permission id
	permissionIDs := make([]uint, 0)

	for _, r := range roleACL {
		if slice.DoesExist(permissionIDs, r.PermissionID) {
			permissionIDs = append(permissionIDs, r.PermissionID)
		}
	}

	if len(permissionIDs) == 0 {
		return nil, nil
	}

	//select * from permissions where id in (?,?,?,?...)
	args := make([]any, len(permissionIDs))

	for i, id := range permissionIDs {
		args[i] = id
	}

	// its error pron area if we had less than one permission id !!!
	query := "select * from permissions where id in (?" + strings.Repeat(",?", len(permissionIDs)-1) +
		")"

	pRows, err := d.conn.Conn().Query(query, args...)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}
	defer pRows.Close()

	permissionTitles := make([]permission.PermissionTitle, 0)

	for pRows.Next() {
		permission, err := scanPermission(pRows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
		}

		permissionTitles = append(permissionTitles, permission.Title)
	}

	if err := pRows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}
	return permissionTitles, nil

}

func scanAccessControl(scanner mysql.Scanner) (accesscontrol.AccessControl, error) {
	const op = "mysql.scanAccessControl"
	var createdAt time.Time
	var acl accesscontrol.AccessControl

	err := scanner.Scan(&acl.ID, &acl.ActorID, &acl.ActorType, &acl.PermissionID, &createdAt)

	return acl, richerror.New(op).
		WithErr(err).WithKind(richerror.KindUnexpected).
		WithMessage(errmsg.ErrorMsgCantScanQueryResult)
}
