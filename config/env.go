package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Loadenv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file!")
	}
}
