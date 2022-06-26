package main

import (
	"encoding/json"
	"fmt"
	global "github.com/go-Server/config"
	"github.com/go-Server/model"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var sessions = map[string]string{}

func SighUp() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var user model.SignupRequestUser
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				global.Logger.Error(err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			//Validate Request Value
			if ValidateRequestUser(user) == false {
				global.Logger.Error("Bad request to Signup")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			//Insert Database
			err = InsertUser(user)
			if err != nil {
				global.Logger.Error(err.Error())
				w.WriteHeader(http.StatusConflict)
				return
			}
			w.WriteHeader(http.StatusCreated)
			global.Logger.Info("생성항 유저 : " + user.Username)
		})
}

func ValidateRequestUser(user model.SignupRequestUser) bool {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return false
	}
	return true
}

func GeneratePassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func DuplicateCheckUser(username string) error {
	var selectedUsername string

	err := global.Db.QueryRow(model.SelectUsernameQuery).Scan(&selectedUsername)
	if err != nil {
		return err
	}
	if selectedUsername == username {
		return fmt.Errorf("이미 사용자가 존재합니다")
	}
	return nil
}

func InsertUser(user model.SignupRequestUser) error {
	//유저 중복검사
	err := DuplicateCheckUser(user.Username)
	if err != nil {
		return err
	}

	//비밀번호 해싱
	hashedPassword, err := GeneratePassword(user.Password)
	if err != nil {
		return err
	}

	//uuid 생성
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	fmt.Println("hash:", hashedPassword)
	fmt.Println("uuid:", uuid)

	_, err = global.Db.Exec(model.InsertUserQuery,
		uuid, user.Username, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}
