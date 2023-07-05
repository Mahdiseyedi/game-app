package main

import (
	"fmt"
	"game-app/entity/user"
	"game-app/repository/mysql"
	"game-app/service/userService"
	"net/http"
)

type UserRegisterHandler struct {
}

func (u UserRegisterHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {

}

func main() {
	test_Register_By_Service()
}

func test_sql() {
	sqlTest := mysql.New()
	srv := userService.Repository(sqlTest)
	srv.RegisterUser(user.User{})
	ph := "13434231"

	res, err := sqlTest.IsPhoneNumberUnique(ph)

	fmt.Println(res)
	fmt.Println(err)
}

func test_Register_By_Service() {
	sqlTest := mysql.New()
	rep := userService.Repository(sqlTest)

	var req userService.RegisterRequest
	req.Name = "hasher"
	req.PhoneNumber = "09541619004"
	req.Password = "mahdi@1234567"
	//srv.RegisterUser(user.User{
	//	ID:          0,
	//	Name:        "mahdi",
	//	PhoneNumber: "0912454545",
	//})

	srv := userService.New(rep)
	resp, qErr := srv.Register(req)

	fmt.Println(resp)
	fmt.Println(qErr)

}

/*
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
*/
