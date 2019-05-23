package main

import "sync"

func customerUsersUpdater(usersLdap []*UserLDAP) {
	customerUsers, _ := getCustomerUsersFromDB()
	newUsers, changedUsers := getCustomerUserUpdates(customerUsers, usersLdap)
	// создаем новых пользователей
	for _, item := range newUsers {
		go createCustomerUser(item)
	}
	// обновляем измененных
	for _, item := range changedUsers {
		go updateCustomerUser(item)
	}
}

func usersUpdater(usersLdap []*UserLDAP) {
	var wg sync.WaitGroup
	// достаем пользователей из бд
	users, _ := getUsersFromDB()
	// получаем обновления и ожидаем создание новых пользователей
	updates := getUserUpdates(users, usersLdap)
	wg.Add(len(updates))
	for _, item := range updates {
		go createUser(item, &wg)
	}
	wg.Wait()
}
