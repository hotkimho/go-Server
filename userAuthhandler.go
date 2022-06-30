package main

import (
	"context"
	"encoding/json"
	"fmt"
	users "github.com/go-Server/auth/users"
	global "github.com/go-Server/config"
	"github.com/go-Server/model"
	"github.com/gofrs/uuid"
	"net/http"
	"time"
)

//일단 DB에 저장해서 구현해보자 사용 X
//var sessions = map[string]string{}

func SighUp() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//Post 요청이 아니면 바로 종료
			if r.Method != http.MethodPost {
				global.Logger.Error("not Post Request")
				http.Error(w, "잘못된 요청입니다", http.StatusBadRequest)
				return
			}

			//json으로 데이터를 가져온다
			var user model.SignupRequestUser
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "잘못된 요청입니다", http.StatusBadRequest)
				return
			}
			fmt.Println(user.Username, user.Password)
			//Validate Request Value
			if users.ValidateRequestUser(user) == false {
				global.Logger.Error("Bad request to Signup")
				http.Error(w, "", http.StatusBadRequest)
				return
			}

			//Insert Database
			err = users.InsertUser(user)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "사용자가 이미 존재합니다", http.StatusConflict)
				return
			}
			http.Error(w, "회원가입이 성공했습니다", http.StatusCreated)
			global.Logger.Info("생성한 유저 : " + user.Username)
		})
}

func Login() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//이미 로그인된 경우는 index 페이지로 이동
			//이동 코드

			//Post 요청이 아닌경우 종료
			if r.Method != http.MethodPost {
				global.Logger.Error("Not Post Request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			//Request Body 데이터를  JSON으로 변환
			var user model.SignupRequestUser
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "잘못된 요청입니다", http.StatusBadRequest)
				return
			}

			err = users.CheckLoginRequest(user.Username, user.Password)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "아이디 또는 비밀번호가 일치하지 않습니다", http.StatusBadRequest)
				return
			}

			newUuid, _ := uuid.NewV4()
			userUuid, ok := users.GetUserUuid(user.Username)
			if ok == false {
				http.Error(w, "아이디 또는 비밀번호가 일치하지 않습니다", http.StatusBadRequest)
				return
			}
			userSession := model.Session{SessionId: newUuid.String(), UserId: userUuid}
			ctx := context.Background()
			err = global.Rdb.Set(ctx, newUuid.String(), userUuid, time.Minute*model.SessionExpiryTime).Err()
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "로그인이 실패했습니다", http.StatusBadRequest)
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "sessionId",
				Value:    userSession.SessionId,
				Secure:   true,
				HttpOnly: true,
				Path:     "/",
				MaxAge:   300,
			})
			userSessiontoJson, _ := json.Marshal(userSession)

			//전역 상태를 관리할 수 없으므로 임시로 세션값을 넘겨줘서 localstorage에 저장하여 사용한다.
			//전역 상태 contextAPI를 배워서 로그인 상태를 저장하고 이 코드는 삭제할것
			http.Error(w, string([]byte(userSessiontoJson)), http.StatusOK)
			fmt.Println("key session id", userSession.SessionId)
			fmt.Println("value user id", userSession.UserId)

			return
		})
}

//로그인을 확인하는 미들웨어
func AuthenticateLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "로그인이 필요합니다.", http.StatusUnauthorized)
				return
			}
			fmt.Println(cookie)
			next.ServeHTTP(w, r)
		},
	)
}
