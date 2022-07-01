package model

import "time"

const PasswordReg = `/^[A-Za-z0-9]{6,12}$/`
const InsertUserQuery string = "INSERT INTO user (uuid, username, password) value(?, ?, ?)"
const InsertSessionQuery string = "INSERT INTO user_session(session_id, user_id) value(?, ?)"
const SelectUsernameQuery string = "SELECT username FROM user WHERE username=?"
const SelectUsernamdAndPassworQuery string = "select username, password from user where username=?"
const SelectUserUuid string = "select uuid from user where username=?"

const SessionExpiryTime time.Duration = 10

type SignupRequestUser struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type User struct {
	Uuid     string
	Username string
	Password string
}

/* user 테이블 생성 쿼리
CREATE TABLE IF NOT EXISTS user (
uuid VARCHAR(36) NOT NULL,
username TEXT NOT NULL,
password TEXT NOT NULL,
created_at DATETIME DEFAULT  CURRENT_TIMESTAMP,
primary key(uuid)
);
*/

/*
Session
삭제된 테이블 및 구조체
생성한 Query
CREATE TABLE IF NOT EXISTS user_session (
session_id TEXT NOT NULL,
user_id VARCHAR(36) NOT NULL,
created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY(user_id)
REFERENCES user(uuid)
);
*/
type Session struct {
	SessionId string
	UserId    string
}
