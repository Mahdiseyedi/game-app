package main

import (
	"encoding/json"
	"fmt"
	"game-app/repository/mysql"
	"game-app/service/authService"
	"game-app/service/userService"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	JwtSignKey           = "jwt_secret"
	AccessTokenSubject   = "at"
	RefreshTokenSubject  = "rt"
	AccessTokenDuration  = time.Hour * 24
	RefreshTokenDuration = time.Hour * 24 * 7
)

func main() {

	http.HandleFunc("/healthcheck", healthCheckHandler)
	http.HandleFunc("/users/register", UserRegisterHandler)
	http.HandleFunc("/users/profile", UserGetProfileHandler)
	http.HandleFunc("/users/login", UserLoginHandler)

	log.Println("server is listening on port 8088...")
	http.ListenAndServe(":8080", nil)
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
	authSrv := authService.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject,
		AccessTokenDuration, RefreshTokenDuration)
	srv := userService.New(mysqlrepo, authSrv)

	res, rErr := srv.Register(rReq)
	if rErr != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, rErr)
		return
	}

	fmt.Fprintf(writer, `{"your user id":"%d"}`, res.User.ID)
}
func UserGetProfileHandler(writer http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		fmt.Fprintf(writer, `{"error":"Invalid Method!"}`)
		return
	}
	authSrv := authService.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject,
		AccessTokenDuration, RefreshTokenDuration)

	authToken := req.Header.Get("Authorization")
	claims, err := authSrv.VerifyToken(authToken)
	if err != nil {
		fmt.Fprintf(writer, `{"error":"Authentication Failed, %s"}`, err)
		return
	}

	mysqlrepo := mysql.New()

	srv := userService.New(mysqlrepo, authSrv)
	res, rErr := srv.GetUserProfile(userService.UserProfileRequest{UserID: claims.UserID})
	if rErr != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, rErr)
		return
	}

	fmt.Fprintf(writer, `{"user profile":"%s"}`, res.Name)
}
func UserLoginHandler(writer http.ResponseWriter, req *http.Request) {
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
	authSrv := authService.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject,
		AccessTokenDuration, RefreshTokenDuration)
	srv := userService.New(mysqlrepo, authSrv)

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
