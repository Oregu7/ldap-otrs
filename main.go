package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// загружаем переменные окружения
	if godotenv.Load() != nil {
		log.Fatal("Error loading .env file")
	}
	users, err := findUsersFromLDAP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(users)

	fmt.Scanln()
}
