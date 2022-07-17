package main

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	global "github.com/hotkimho/go-Server/config"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"time"
)

func GetMysqlConnector() *sql.DB {
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASSWORD")
	address := os.Getenv("DBADDRESS")
	dbName := os.Getenv("DBNAME")

	mysqlConfig := mysql.Config{
		User:                 user,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 address,
		Collation:            "utf8mb4_general_ci",
		Loc:                  time.UTC,
		MaxAllowedPacket:     4 << 20.,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
		DBName:               dbName,
	}
	connector, err := mysql.NewConnector(&mysqlConfig)
	if err != nil {
		panic(err)
	}
	db := sql.OpenDB(connector)
	return db
}

func indexHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			http.Error(w, "로그인이 성공했습니다", http.StatusOK)
			//fmt.Fprintf(w, users.Username)
		},
	)
}

func main() {
	//환경 변수 설정
	config := godotenv.Load(".env")
	if config != nil {
		log.Panic("Error loading dotenv")
	}
	//데이터 베이스 연결 객체 설정
	global.Db = GetMysqlConnector()

	//로거 설정
	global.Logger, _ = zap.NewDevelopment()
	defer global.Logger.Sync()
	defer global.Db.Close()

	//redis test

	global.Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	router := mux.NewRouter()
	router.Handle("/", indexHandler()).Methods("GET", "POST")
	router.Handle("/auth/signup", SighUp()).Methods("POST")
	router.Handle("/auth/login", Login()).Methods("POST")
	router.Handle("/auth/logout", SessionAuthenticate(Logout())).Methods("GET")

	router.Handle("/board", GetBoard()).Methods("GET")
	router.Handle("/board/post", GetPageOfBoard()).Methods("GET")
	router.Handle("/board/post", SessionAuthenticate(CreatePost())).Methods("POST")
	router.Handle("/board/post", SessionAuthenticate(DeletePost())).Methods("DELETE")
	router.Handle("/board/post", SessionAuthenticate(EditPost())).Methods("PATCH")

	fmt.Println("서버가 시작되었습니다. ")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://www.hotkimho.com"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch},
	})

	handler := corsHandler.Handler(router)
	_ = http.ListenAndServe("127.0.0.1:8000", handler)

}
