package users

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	global "github.com/hotkimho/go-Server/config"
	"github.com/hotkimho/go-Server/model"
	"golang.org/x/crypto/bcrypt"
)

func GetUserUuid(username string) (string, bool) {
	var selectedUuid string

	err := global.Db.QueryRow(model.SelectUserUuid, username).Scan(&selectedUuid)
	if err != nil {
		return "", false
	}
	return selectedUuid, true
}

func CheckLoginRequest(username, password string) error {
	var (
		selectedUsername string
		selectedPassword string
	)

	//회원이 있는지 검사
	err := DuplicateCheckUser(username)
	if err == nil {
		//프론트에 전달하는 로그인 에러는
		//1. 사용자가 없습니다
		//2. 사용자 또는 비밀번호가 일치하지 않습니다.
		//보안을 위해 2번을 해야하는가
		return fmt.Errorf("아이디 또는 비밀번호가 일치하지 않습니다")
	}

	//이미 중복검사를 통해 유저가 있는걸 확인했으니 유저가 없는 경우는 확인 X
	resultQuery := global.Db.QueryRow(model.SelectUsernamdAndPassworQuery, username)
	if resultQuery.Err() != nil {
		return resultQuery.Err()
	}
	_ = resultQuery.Scan(&selectedUsername, &selectedPassword)

	//비밀번호 검사
	err = bcrypt.CompareHashAndPassword([]byte(selectedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func GetUuidInSession(sessionId string) (username string, err error) {
	ctx := context.Background()
	username, err = global.Rdb.Get(ctx, sessionId).Result()
	//redis를 조회하며, 사용자가 있으면 username을 리턴
	if err == redis.Nil {
		return "", err
	} else if err != nil {
		return "", err
	}
	return username, nil
}
