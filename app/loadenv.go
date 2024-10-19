package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("No .env file found. Using system environment variables.")
		return
	}
	fmt.Println(".env loaded successfully.")
}
