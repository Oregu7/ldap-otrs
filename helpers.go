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

// getUpdates получаем список новых пользователей из LDAP
func getUpdates(users []*User, usersLdap []*UserLDAP) []*UserLDAP {
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
