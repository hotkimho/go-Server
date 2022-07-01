package model

type Post struct {
	PostId   int    `json:"id"`
	Title    string `json:"title"`
	Writer   string `json:"writer"`
	Content  string `json:"content"`
	CreateAt string `json:"created_at"`
	View     int    `json:"view"`
}

type RequestPost struct {
	Title   string
	Content string
}

type ReponsePost struct {
	Title   string `json:"title"`
	Writer  string `json:"writer"`
	Content string `json:"content"`
}

/* post 테이블 생성 쿼리
CREATE TABLE IF NOT EXISTS post (
postId int auto_increment,
title TEXT NOT NULL,
writer TEXT NOT NULL,
content TEXT,
view int NOT NULL DEFAULT 0,
created_at DATE DEFAULT (current_date),
userId VARCHAR(36) NOT NULL,
primary key(postId),
foreign key (userId) references user(uuid)
);

*/
