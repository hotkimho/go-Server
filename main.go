package main

import (
	"context"
	"database/sql"
	"fmt"
	global "github.com/go-Server/config"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
			fmt.Println("index page")
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
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	err := rdb.Set(ctx, "kim", "ho", 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rdb.Get(ctx, "kim").Result())

	router := mux.NewRouter()
	router.Handle("/", indexHandler()).Methods("POST")
	router.Handle("/auth/signup", SighUp()).Methods("POST")
	router.Handle("/auth/login", Login()).Methods("POST")
	router.Handle("/auth/test", LoginAuth(indexHandler())).Methods("POST")
	_ = http.ListenAndServe("127.0.0.1:3000", router)
}
