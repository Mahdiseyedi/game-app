package main

import (
	"encoding/json"
	"fmt"
	"game-app/repository/mysql"
	"game-app/service/userService"
	"io"
	"log"
	"net/http"
)

type UserRegisterHandler struct {
}

func (u UserRegisterHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {

}

func main() {

	http.HandleFunc("/users/register", userRegisterHandler)
	log.Println("server is listening on port 8088...")
	http.ListenAndServe(":8080", nil)
}

func test_Register_By_Service() {
	sqlTest := mysql.New()
	rep := userService.Repository(sqlTest)

	var req userService.UserProfileRequest

	req.UserID = 7

	//srv.RegisterUser(user.User{
	//	ID:          0,
	//	Name:        "mahdi",
	//	PhoneNumber: "0912454545",
	//})

	srv := userService.New(rep)
	resp, qErr := srv.GetUserProfile(req)

	fmt.Println(resp)
	fmt.Println(qErr)

}

func userRegisterHandler(writer http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		fmt.Fprintf(writer, `{"error":"Invalid Method!"}`)
		return
	}

	var rReq userService.RegisterRequest
	body, _ := io.ReadAll(req.Body)
	err := json.Unmarshal(body, &rReq)
	if err != nil {
		fmt.Fprintf(writer, `{"error":"invalid json format..."}`)
		return
	}

	mysqlrepo := mysql.New()
	srv := userService.New(mysqlrepo)

	res, rErr := srv.Register(rReq)
	if rErr != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, rErr)
		return
	}

	fmt.Fprintf(writer, `{"your user id":"%d"}`, res.User.ID)
}

func healthCheckHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, `{"message": "everything is good!"}`)
}
