package main

import (
	"encoding/json"
	"fmt"
	"game-app/entity"
	"game-app/repository/mysql"
	"game-app/service/userService"
	"io"
	"net/http"
)

type UserRegisterHandler struct {
}

func (u UserRegisterHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {

}

func main() {

	http.HandleFunc("/users/register", func(writer http.ResponseWriter, request *http.Request) {

		if request.Method == http.MethodGet {
			fmt.Println("Get")
			data := map[string]string{"message": "hi\n", "name": "mahdi"}
			js, err := json.Marshal(data)
			if err != nil {
				fmt.Println(err)
			}
			writer.Header().Set("Content-Type", "application/json")
			writer.Write(js)
		}
		if request.Method == http.MethodPost {
			var t userService.RegisterRequest
			fmt.Println("post")
			u, _ := io.ReadAll(request.Body)
			err := json.Unmarshal(u, &t)
			if err != nil {
				fmt.Println("this is not valid json format", err)
				return
			}

			mysqlrepo := mysql.New()
			th := userService.New(mysqlrepo)

			fmt.Println(t)

			k, e := th.Register(t)
			if e != nil {
				fmt.Println(e)
			}
			fmt.Println(k)
		}

	})

	http.ListenAndServe(":8080", nil)
}

func test_sql() {
	sqlTest := mysql.New()
	srv := userService.Repository(sqlTest)
	srv.RegisterUser(entity.User{})
	ph := "13434231"

	res, err := sqlTest.IsPhoneNumberUnique(ph)

	fmt.Println(res)
	fmt.Println(err)
}
