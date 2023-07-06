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

const (
	JwtSignKey = "jwt_secret"
)

func main() {

	http.HandleFunc("/healthcheck", healthCheckHandler)
	http.HandleFunc("/users/register", UserRegisterHandler)

	//TODO - implementing GetUserProfile handler
	http.HandleFunc("/users/profile", UserGetProfileHandler)

	//TODO - implementing Login handler
	http.HandleFunc("/users/login", UserLoginHandler)

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

	srv := userService.New(rep, JwtSignKey)
	resp, qErr := srv.GetUserProfile(req)

	fmt.Println(resp)
	fmt.Println(qErr)

}

func UserRegisterHandler(writer http.ResponseWriter, req *http.Request) {

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
	srv := userService.New(mysqlrepo, JwtSignKey)

	res, rErr := srv.Register(rReq)
	if rErr != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, rErr)
		return
	}

	fmt.Fprintf(writer, `{"your user id":"%d"}`, res.User.ID)
}
func UserGetProfileHandler(writer http.ResponseWriter, req *http.Request) {

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
	srv := userService.New(mysqlrepo, JwtSignKey)

	res, rErr := srv.Register(rReq)
	if rErr != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, rErr)
		return
	}

	fmt.Fprintf(writer, `{"your user id":"%d"}`, res.User.ID)
}
func UserLoginHandler(writer http.ResponseWriter, req *http.Request) {

	r, e := userService.CreateToken(12, JwtSignKey)
	if e != nil {
		fmt.Fprintf(writer, `{"error": %s `, e)
		return
	}
	fmt.Println(r)

	if req.Method != http.MethodPost {
		fmt.Fprintf(writer, `{"error":"Invalid Method!"}`)
		return
	}

	var lReq userService.LoginRequest
	body, _ := io.ReadAll(req.Body)
	err := json.Unmarshal(body, &lReq)
	if err != nil {
		fmt.Fprintf(writer, `{"error":"invalid json format..."}`)
		return
	}

	mysqlrepo := mysql.New()
	srv := userService.New(mysqlrepo, JwtSignKey)

	res, rErr := srv.Login(lReq)
	if rErr != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, rErr)
		return
	}

	loginResult, mErr := json.Marshal(res)
	if mErr != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, mErr)
		return
	}

	writer.Write(loginResult)
}

func healthCheckHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, `{"message": "everything is good!"}`)
}
