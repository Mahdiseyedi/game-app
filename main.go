package main

import (
	"encoding/json"
	"fmt"
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
			ur := map[string]string{}
			fmt.Println("post")
			u, _ := io.ReadAll(request.Body)
			err := json.Unmarshal(u, &ur)
			if err != nil {
				fmt.Println("this is not valid json format", err)
				return
			}

			var t userService.RegisterRequest
			t.Name = ur["name"]
			t.PhoneNumber = ur["phoneNumber"]
			fmt.Println(t.Name)
			fmt.Println(t.PhoneNumber)
			var s userService.Service
			iu, erq := s.Register(t)
			if erq != nil {
				fmt.Println("smth wrong in Db, ", erq)
				return
			}
			fmt.Println(iu)

		}

	})

	http.ListenAndServe(":8080", nil)
}

func test_sql() {
	sqlTest := mysql.New()

	ph := "13434231"

	res, err := sqlTest.IsPhoneNumberUnique(ph)

	fmt.Println(res)
	fmt.Println(err)
}
