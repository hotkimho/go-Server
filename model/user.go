package model

const PasswordReg = `/^[A-Za-z0-9]{6,12}$/`

type SignupRequestUser struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type User struct {
	Uuid     string
	Username string
	Password string
}

// Session
// 생성한 Query
//CREATE TABLE IF NOT EXISTS user_session (
//session_id TEXT NOT NULL,
//user_id VARCHAR(36) NOT NULL,
//FOREIGN KEY(user_id)
//REFERENCES user(uuid)
//);
///*
type Session struct {
	SessionId string
	UserId    string
}
