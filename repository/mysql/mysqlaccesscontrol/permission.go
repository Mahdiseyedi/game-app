package mysqlaccesscontrol

import (
	"game-app/entity/permission"
	"game-app/repository/mysql"
	"time"
)

func scanPermission(scanner mysql.Scanner) (permission.Permission, error) {
	var createdAt time.Time
	var p permission.Permission

	err := scanner.Scan(&p.ID, &p.Title, &createdAt)

	return p, err
}
