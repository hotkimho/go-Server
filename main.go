package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func createTable(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS
	
	
`
}

func main() {
	config := godotenv.Load(".env")
	if config != nil {
		log.Fatal("Error loading dotenv")
	}

	db, err := sql.Open("mysql", os.Getenv("DBURL"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ok", db.Ping())
}
