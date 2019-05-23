package main

import (
	"log"
	"sync"
)

// updateUsersTask синхронизация агентов OTRS и пользователей из LDAP
func updateUsersTask() {
	var wg sync.WaitGroup
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
