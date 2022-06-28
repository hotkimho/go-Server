package model

const PasswordReg = `/^[A-Za-z0-9]{6,12}$/`
const InsertUserQuery string = "INSERT INTO user (uuid, username, password) value(?, ?, ?)"
const InsertSessionQuery string = "INSERT INTO user_session(session_id, user_id) value(?, ?)"
const SelectUsernameQuery string = "SELECT username FROM user WHERE username=?"
const SelectUsernamdAndPassworQuery string = "select username, password from user where username=?"
const SelectUserUuid string = "select uuid from user where username=?"

type SignupRequestUser struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type User struct {
	Uuid     string
	Username string
	Password string
}

/*
CREATE TABLE IF NOT EXISTS users (
uuid VARCHAR(36) NOT NULL,
username TEXT NOT NULL,
password TEXT NOT NULL,
created_at DATETIME DEFAULT  CURRENT_TIMESTAMP,
primary key(uuid)
);
*/

// Session
// 생성한 Query
//CREATE TABLE IF NOT EXISTS user_session (
//session_id TEXT NOT NULL,
//user_id VARCHAR(36) NOT NULL,
//created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
//FOREIGN KEY(user_id)
//REFERENCES user(uuid)
//);
///*
type Session struct {
	SessionId string
	UserId    string
}

/*
CREATE EVENT
	IF NOT EXISTS DELETE_LOGIN_SESSION
ON SCHEDULE
	EVERY 15 MINUTE


*/

/*
DELIMITER $$
CREATE PROCEDURE DELETE_LOGIN_SESSION()
BEGIN

*/
