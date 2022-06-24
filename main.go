package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-Server/model"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var db *sql.DB
var logger *zap.Logger

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
			var user model.SignupRequestUser
			json.NewDecoder(r.Body).Decode(&user)

			body, _ := ioutil.ReadAll(r.Body)

			fmt.Println(user, body)

			fmt.Println("12")
		},
	)
}

//
//func sighUp() http.Handler {
//	return http.HandlerFunc(
//		func(w http.ResponseWriter, r *http.Request) {
//			fmt.Println("123")
//			w.Write([]byte("hi"))
//			body := r.Body
//			fmt.Println(body)
//		})
//}

func main() {
	//환경 변수 설정
	config := godotenv.Load(".env")
	if config != nil {
		log.Panic("Error loading dotenv")
	}
	//데이터 베이스 연결 객체 설정
	db = GetMysqlConnector()
	//로거 설정
	logger, _ = zap.NewDevelopment()
	defer logger.Sync()

	router := mux.NewRouter()
	router.Handle("/", indexHandler()).Methods("POST")
	router.Handle("/auth/signup", sighUp()).Methods("POST")
	fmt.Println(uuid.New())
	_ = http.ListenAndServe("127.0.0.1:3000", router)
}
