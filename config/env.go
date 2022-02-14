package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Loadenv() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println("Error loading .env file!")
		}
	}

}
