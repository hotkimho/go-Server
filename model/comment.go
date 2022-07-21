package model

const InsertComment string = "INSERT INTO comment (writer, content, userId, postId) value(?, ?, ?, ?)"
const SelectCommentInpostId = "SELECT writer, content, created_at from comment where postId=?"

type Comment struct {
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
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
foreign key (postId) references post(postId)
);
*/
