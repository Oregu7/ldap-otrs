package main

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	var wg sync.WaitGroup
	// загружаем переменные окружения
	if godotenv.Load() != nil {
		log.Fatal("Error loading .env file")
	}
	// достаем пользователей из бд
	users, err := getUsersFromDB()
	if err != nil {
		log.Fatal(err)
	}
	// достаем пользователей из ldap
	usersLdap, err := findUsersFromLDAP()
	if err != nil {
		log.Fatal(err)
	}
	// получаем обновления и ожидаем создание новых пользователей
	updates := getUserUpdates(users, usersLdap)
	wg.Add(len(updates))
	for _, item := range updates {
		go createUser(item, &wg)
	}
	wg.Wait()
	log.Println("[Синхронизация завершена...]")
}
