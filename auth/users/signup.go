package users

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	global "github.com/hotkimho/go-Server/config"
	"github.com/hotkimho/go-Server/model"
	"golang.org/x/crypto/bcrypt"
)

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

	err := global.Db.QueryRow(model.SelectUsernameQuery, username).Scan(&selectedUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
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
	newUuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	fmt.Println("hash:", hashedPassword)
	fmt.Println("uuid:", newUuid)

	_, err = global.Db.Exec(model.InsertUserQuery,
		newUuid.String(), user.Username, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}
