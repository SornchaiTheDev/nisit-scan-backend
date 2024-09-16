package libs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func InitEnv() {
	isProd, err := strconv.ParseBool(os.Getenv("IS_PROD"))
	if err != nil {
		isProd = false
	}

	if isProd {
		return
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Loaded .env file")
}
