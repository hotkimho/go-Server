package model

const InsertComment string = "INSERT INTO comment (writer, content, userId, postId) value(?, ?, ?, ?)"

type RequestComment struct {
	PostId  string `json:"postId"`
	Content string `json:"content"`
}

/*
CREATE TABLE IF NOT EXISTS comment (
commentId int auto_increment,
writer TEXT NOT NULL,
content TEXT,
created_at DATE DEFAULT (current_date),
userId VARCHAR(36) NOT NULL,
postId int NOT NULL,
primary key(commentId),
foreign key (userId) references user(uuid),
foreign key (commentId) references post(postId)
);
*/
