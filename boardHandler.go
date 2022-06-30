package main

import (
	"encoding/json"
	"fmt"
	global "github.com/go-Server/config"
	"github.com/go-Server/model"
	"net/http"
	"strconv"
)

//페이지 번호에 맞는 게시글을 가져옴
func GetBoard() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			const limit int = 10
			getPageOfString := r.URL.Query().Get("page")
			page, _ := strconv.Atoi(getPageOfString)
			if page > 0 {
				page = page - 1
			}
			fmt.Println(page)
			fmt.Println(page*10, page*10+limit)
			selectAllBoardQuery := fmt.Sprintf("SELECT postId, title, writer, content , created_at, view FROM post limit %d, %d", page*10, page*10+limit)
			fmt.Println(selectAllBoardQuery)
			rows, err := global.Db.Query(selectAllBoardQuery)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "게시글을 불러오지 못했습니다", http.StatusNotFound)
				return
			}
			defer rows.Close()

			var posts []model.Post
			for rows.Next() {
				var post model.Post
				_ = rows.Scan(&post.PostId, &post.Title, &post.Writer,
					&post.Content, &post.CreateAt, &post.View)
				posts = append(posts, post)
			}
			err = rows.Err()
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "게시글을 불러오지 못했습니다", http.StatusNotFound)
				return
			}

			postToJson, err := json.Marshal(posts)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "게시글을 불러오지 못했습니다", http.StatusNotFound)
				return
			}

			http.Error(w, "", http.StatusOK)
			w.Write(postToJson)
		},
	)
}

func CreatePost() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				global.Logger.Error("Not Post Request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			//Request Body 데이터를  JSON으로 변환
			var createPost model.RequestPost
			err := json.NewDecoder(r.Body).Decode(&createPost)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "잘못된 요청입니다", http.StatusBadRequest)
				return
			}
			cookie, _ := r.Cookie("sessionId")
			fmt.Println(cookie.Value)
		},
	)
}
