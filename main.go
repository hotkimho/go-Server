package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	config := godotenv.Load(".env")
	if config != nil {
		log.Fatal("Error loading dotenv")
	}
	fmt.Println(os.Getenv("TEST"))
}
