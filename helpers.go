package main

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
)

// getMD5Hash возвращает хеш строки
func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// createUsersPattern создаем паттерн из логинов пользователей для поиска
func createUsersPattern(users []*User) string {
	var pattern string
	for i, user := range users {
		pattern += getMD5Hash(user.Login)
		if len(users)-1 != i {
			pattern += "|"
		}
	}

	return pattern
}

func createUsersHashMap(users []*CustomerUser) map[string]string {
	hashMap := make(map[string]string)
	for _, user := range users {
		hashMap[user.getHashMapKey()] = user.getPropsToken()
	}
	return hashMap
}

// getUpdates получаем список новых пользователей из LDAP
func getUserUpdates(users []*User, usersLdap []*UserLDAP) []*UserLDAP {
	updates := []*UserLDAP{}
	pattern := createUsersPattern(users)
	for _, item := range usersLdap {
		// ищем обновления
		loginHash := getMD5Hash(item.Login)
		mtch, _ := regexp.MatchString(pattern, loginHash)
		if len(pattern) == 0 || !mtch {
			updates = append(updates, item)
		}
	}

	return updates
}

// getUpdates получаем список новых и changed клиентов из LDAP
func getCustomerUserUpdates(customerUsers []*CustomerUser, usersLdap []*UserLDAP) ([]*UserLDAP, []*UserLDAP) {
	updates := []*UserLDAP{}
	changed := []*UserLDAP{}

	hashMap := createUsersHashMap(customerUsers)
	for _, item := range usersLdap {
		// ищем обновления
		if token, ok := hashMap[item.getHashMapKey()]; ok {
			if token != item.getPropsToken() {
				changed = append(changed, item)
			}
		} else {
			updates = append(updates, item)
		}
	}

	return updates, changed
}
