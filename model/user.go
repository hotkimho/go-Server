package model

type SignupRequestUser struct {
	Username string
	Password string
}

// Session
// 생성한 Query
// CREATE TABLE IF NOT EXISTS user_session (
// session_id TEXT NOT NULL,
// user_id binary(16) NOT NULL,
// FOREIGN KEY(user_id)
// REFERENCES user(uuid)
// );
///*
type Session struct {
	SessionId string
	UserId    string
}
