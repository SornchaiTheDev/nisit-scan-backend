package libs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func InitEnv() {
	isDocker, err := strconv.ParseBool(os.Getenv("IS_PROD"))
	if err != nil {
		isDocker = false
	}

	if isDocker {
		return
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Loaded .env file")
}
