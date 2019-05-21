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
	users, err := getUsersFromDB()
	if err != nil {
		log.Fatal(err)
	}

	updates := getUpdates(users, []*UserLDAP{&UserLDAP{Username: "roo@localhost"}})
	fmt.Println(updates)
}
