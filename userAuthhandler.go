package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	users "github.com/hotkimho/go-Server/auth/users"
	global "github.com/hotkimho/go-Server/config"
	"github.com/hotkimho/go-Server/model"
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
				MaxAge:   600,
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
func SessionAuthenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//쿠키에 있는 세션값을 조회
			cookie, err := r.Cookie("sessionId")
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "로그인이 필요합니다.", http.StatusUnauthorized)
				return
			}

			//세션의 쿠키값이 있으면 사용자로 인식
			//로그인 확인할 땐, 레디스에서 확인하고 DB에 쿼리를 날리지 않는다!!!!!
			ctx := context.Background()
			_, err = global.Rdb.Get(ctx, cookie.Value).Result()
			fmt.Println(err)
			if err == redis.Nil {
				global.Logger.Error(err.Error())
				http.Error(w, "로그인이 필요합니다.", http.StatusUnauthorized)
				return
			} else if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "접근에 실패했습니다.", http.StatusUnauthorized)
				return
			}

			//이 과정을 거치면 사용자로 인식하고 다음 미들웨어, 핸들러로 이동
			next.ServeHTTP(w, r)
		},
	)
}

func Logout() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("start logout")
			cookie, err := r.Cookie("sessionId")
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "로그인이 필요합니다.", http.StatusUnauthorized)
				return
			}

			fmt.Println(cookie)
			ctx := context.Background()
			err = global.Rdb.Del(ctx, cookie.Value).Err()
			if err == redis.Nil {
				global.Logger.Error(err.Error())
				http.Error(w, "로그인이 필요합니다.", http.StatusUnauthorized)
				return
			} else if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "로그아웃이 실패했습니다.", http.StatusNotFound)
				return
			}
			removeCookie := http.Cookie{
				Name:     "sessionId",
				Value:    "",
				MaxAge:   0,
				Path:     "/",
				HttpOnly: true,
			}
			http.SetCookie(w, &removeCookie)
			http.Error(w, "로그아웃이 완료되었습니다.", http.StatusOK)
		})
}
