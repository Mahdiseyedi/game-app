package mysqlaccesscontrol

import (
	"game-app/entity/permission"
	"game-app/pkg/errmsg"
	"game-app/pkg/richerror"
	"game-app/repository/mysql"
	"time"
)

func scanPermission(scanner mysql.Scanner) (permission.Permission, error) {
	const op = "mysqlaccesscontrol.scanPermission"
	var createdAt time.Time
	var p permission.Permission

	err := scanner.Scan(&p.ID, &p.Title, &createdAt)

	return p, richerror.New(op).WithErr(err).
		WithMessage(errmsg.ErrorMsgCantScanQueryResult).
		WithKind(richerror.KindNotAcceptable)
}
