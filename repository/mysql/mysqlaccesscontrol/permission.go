package mysqlaccesscontrol

import (
	"game-app/entity/permission"
	"game-app/repository/mysql"
)

func scanPermission(scanner mysql.Scanner) (permission.Permission, error) {
	var createdAt []uint8
	var p permission.Permission

	err := scanner.Scan(&p.ID, &p.Title, &createdAt)

	return p, err
}
