package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	global "github.com/hotkimho/go-Server/config"
	"github.com/hotkimho/go-Server/model"
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

			selectAllBoardQuery := fmt.Sprintf("SELECT postId, title, writer, content , created_at, view FROM post limit %d, %d", page*10, page*10+limit)
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
			//Request Body 데이터를  JSON으로 변환
			var createPost model.RequestPost
			err := json.NewDecoder(r.Body).Decode(&createPost)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "잘못된 요청입니다", http.StatusBadRequest)
				return
			}

			cookie, err := r.Cookie("sessionId")
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "로그인이 필요합니다.", http.StatusUnauthorized)
				return
			}

			ctx := context.Background()
			getUuidFromCookie, err := global.Rdb.Get(ctx, cookie.Value).Result()
			if err == redis.Nil {
				global.Logger.Error(err.Error())
				http.Error(w, "로그인이 필요합니다.", http.StatusUnauthorized)
				return
			} else if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "글쓰기가 실패했습니다..", http.StatusUnauthorized)
				return
			}

			//쿠키의 세션으로 사용자의 닉네임을 가져온다.
			var username string
			err = global.Db.QueryRow("SELECT username FROM user WHERE uuid=?", getUuidFromCookie).Scan(&username)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "글쓰기가 실패했습니다..", http.StatusUnauthorized)
				return
			}

			//가져온 사용자 정보로 글을 생성한다.
			_, err = global.Db.Exec("INSERT INTO post(title, content, writer, userId) VALUES (?, ?, ?, ?)",
				createPost.Title, createPost.Content, username, getUuidFromCookie)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "글쓰기가 실패했습니다..", http.StatusUnauthorized)
				return
			}
			http.Error(w, "글쓰기가 성공했습니다", http.StatusCreated)
			return
		},
	)
}

func DeletePost() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			postId := r.URL.Query().Get("postId")
			cookie, err := r.Cookie("sessionId")
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "로그인이 필요합니다.", http.StatusUnauthorized)
				return
			}

			//삭제할 게시글이 사용자의 소유글인지 확인
			ctx := context.Background()
			getUuidFromCookie, err := global.Rdb.Get(ctx, cookie.Value).Result()
			if err == redis.Nil {
				global.Logger.Error(err.Error())
				http.Error(w, "로그인이 필요합니다.", http.StatusUnauthorized)
				return
			} else if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "글쓰기가 실패했습니다..", http.StatusUnauthorized)
				return
			}

			//삭제할 게시글의 소유자를 가져온다.
			var userUuid string
			err = global.Db.QueryRow("SELECT userId FROM post WHERE postId=?", postId).Scan(&userUuid)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "글쓰기가 실패했습니다..", http.StatusUnauthorized)
				return
			}
			//삭제하는 사람과 삭제할 사람이 동일한 사람인지 검사
			if getUuidFromCookie != userUuid {
				global.Logger.Error(err.Error())
				http.Error(w, "글쓰기가 실패했습니다..", http.StatusUnauthorized)
				return
			}

			//가져온 사용자 정보로 글을 생성한다.
			_, err = global.Db.Exec("DELETE FROM post WHERE postId = ? ", postId)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "글쓰기가 실패했습니다..", http.StatusUnauthorized)
				return
			}
			http.Error(w, "게시글 삭제가 성공했습니다", http.StatusOK)
			return
		},
	)
}

func EditPost() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//수정할 게시글의 ID
			postId := r.URL.Query().Get("postId")

			//Request Body 데이터를  JSON으로 변환
			var editPost model.RequestPost
			err := json.NewDecoder(r.Body).Decode(&editPost)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "잘못된 요청입니다", http.StatusBadRequest)
				return
			}

			//쿠키에 담긴 세션ID를 가져온다.
			cookie, err := r.Cookie("sessionId")
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "로그인이 필요합니다.", http.StatusUnauthorized)
				return
			}

			//레디스 서버에 uuid를 얻기위해 세션ID를 키값 사용합니다.
			ctx := context.Background()
			getUuidFromCookie, err := global.Rdb.Get(ctx, cookie.Value).Result()
			if err == redis.Nil {
				global.Logger.Error(err.Error())
				http.Error(w, "로그인이 필요합니다.", http.StatusUnauthorized)
				return
			} else if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "글수정이 실패했습니다.", http.StatusBadRequest)
				return
			}

			var userUuid string
			err = global.Db.QueryRow("SELECT userId FROM post WHERE postId=?", postId).Scan(&userUuid)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "글수정이 실패했습니다.", http.StatusBadRequest)
				return
			}
			//삭제하는 사람과 삭제할 사람이 동일한 사람인지 검사
			if getUuidFromCookie != userUuid {
				global.Logger.Error(err.Error())
				http.Error(w, "글수정이 실패했습니다.", http.StatusBadRequest)
				return
			}

			//가져온 사용자 정보로 글을 생성한다.
			_, err = global.Db.Exec("UPDATE post set title = ?, content = ? WHERE postId = ?",
				editPost.Title, editPost.Content, postId)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "글수정이 실패했습니다..", http.StatusBadRequest)
				return
			}
			http.Error(w, "글 수정이 성공했습니다", http.StatusOK)
			return
		},
	)
}

func GetPageOfBoard() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//수정할 게시글의 ID
			postId := r.URL.Query().Get("postId")

			//해당ㅇ 게시글에 맞는 글의 정보를 가져온다.
			var resPost model.ReponsePost
			err := global.Db.QueryRow("SELECT title, writer, content FROM post WHERE postId=?", postId).Scan(&resPost.Title, &resPost.Writer, &resPost.Content)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "게시글 불러오기가 실패했습니다.", http.StatusNotFound)
				return
			}

			postToJson, err := json.Marshal(resPost)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "게시글 불러오기가 실패했습니다.", http.StatusNotFound)
				return
			}

			http.Error(w, "", http.StatusOK)
			w.Write(postToJson)
			return
		},
	)
}

func CreateComment() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var requestComment model.RequestComment
			err := json.NewDecoder(r.Body).Decode(&requestComment)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "잘못된 요청입니다", http.StatusBadRequest)
				return
			}
			postId, err := strconv.Atoi(requestComment.PostId)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "잘못된 게시글입니다.", http.StatusBadRequest)
				return
			}

			//이미 로그인 검증을 했으므로 쿠키에선 에러처리를 하지 않는다.
			cookie, _ := r.Cookie("sessionId")
			uuid, err := getUuidInSession(cookie.Value)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "사용자 인증조회에 실패했습니다. 다시 시도해주세요", http.StatusUnauthorized)
				return
			}
			//uuid를 통해 username을 가져온다.

			var username string
			err = global.Db.QueryRow(model.SelectUsernameInUuid, uuid).Scan(&username)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "사용자 조회에 실패했습니다.", http.StatusUnauthorized)
				return
			}
			//검증이 끝난후, 입력받은 댓글을 데이터베이스에 저장한다.
			_, err = global.Db.Exec(model.InsertComment, username, requestComment.Content, uuid, postId)
			if err != nil {
				global.Logger.Error(err.Error())
				http.Error(w, "댓글 저장에 실패했습니다.", http.StatusBadRequest)
				return
			}

			http.Error(w, "댓글쓰기가 성공했습니다.", http.StatusOK)
		})
}

func getUuidInSession(sessionId string) (username string, err error) {
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
