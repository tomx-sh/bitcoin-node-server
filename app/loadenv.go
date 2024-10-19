package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("ENV")

	if env == "production" {
		err := godotenv.Load(".env.production")
		if err != nil {
			log.Fatal("Error loading .env.production file")
		}

	} else {
		if env == "" {
			os.Setenv("ENV", "local")
		}

		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatal("Error loading .env.local file")
		}
	}

	fmt.Println("Environment:", os.Getenv("ENV"))
}
