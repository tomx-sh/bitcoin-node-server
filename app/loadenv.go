package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("No .env file found. Using system environment variables.")

	} else {
		fmt.Println(".env loaded successfully.")
	}

	fmt.Println("RPC_ENV:", os.Getenv("RPC_ENV"))
	fmt.Println("RPC_USERNAME:", os.Getenv("RPC_USERNAME"))
	fmt.Println("RPC_PASSWORD:", os.Getenv("RPC_PASSWORD"))
	fmt.Println("RPC_URL:", os.Getenv("RPC_URL"))
	fmt.Println("RPC_API_KEY:", os.Getenv("RPC_API_KEY"))
	fmt.Println("SSL_CERTIFICATES_PATH:", os.Getenv("SSL_CERTIFICATES_PATH"))
}
