package main

import (
	"fmt"
	"game-app/repository/mysql"
)

func main() {
	sqlTest := mysql.New()

	ph := "13434"

	res, _ := sqlTest.IsPhoneNumberUnique(ph)

	fmt.Println(res)

}
