package main

import (
	"crypto/tls"
	"fmt"
	"os"

	ldap "gopkg.in/ldap.v3"
)

// UserLDAP данные пользователя
type UserLDAP struct {
	FullName, LastName, FirstName, Login string
	Company, Mail, Phone                 string
	Password                             string
}

// findUsersFromLDAP достаем пользователей из LDAP
func findUsersFromLDAP() ([]*UserLDAP, error) {
	users := []*UserLDAP{}
	// подключаемся к серверу ldap
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	l, err := ldap.DialTLS("tcp", os.Getenv("LDAP_HOST"), tlsConfig)
	defer l.Close()
	if err != nil {
		return users, err
	}
	// подключаем пользователя
	err = l.Bind(os.Getenv("LDAP_LOGIN"), os.Getenv("LDAP_PASSWORD"))
	if err != nil {
		return users, err
	}
	// ищем данные
	result, err := l.Search(&ldap.SearchRequest{
		BaseDN:       os.Getenv("LDAP_BASEDN"),
		Scope:        2,
		DerefAliases: 0,
		SizeLimit:    300,
		Filter:       "(userPrincipalName=*)",
		Attributes:   []string{"mail", "sn", "givenName", "userPrincipalName", "displayName", "company", "mail", "telephoneNumber"},
		Controls:     nil,
	})
	if err != nil {
		return users, err
	}

	for _, entry := range result.Entries {
		// пропускаем без фамилии
		if len(entry.GetAttributeValue("sn")) == 0 {
			continue
		}
		user := &UserLDAP{
			entry.GetAttributeValue("displayName"),
			entry.GetAttributeValue("sn"),
			entry.GetAttributeValue("givenName"),
			entry.GetAttributeValue("userPrincipalName"),
			entry.GetAttributeValue("company"),
			entry.GetAttributeValue("mail"),
			entry.GetAttributeValue("telephoneNumber"),
			"daaad6e5604e8e17bd9f108d91e26afe6281dac8fda0091040a7a6d7bd9b43b5",
		}
		users = append(users, user)
	}

	return users, nil
}

func (user *UserLDAP) successLog() {
	fmt.Printf("[OK] Пользователь: (%s %s)\n", user.FullName, user.Login)
}

func (user *UserLDAP) errorLog(err error) {
	fmt.Printf("[ERROR] Пользователь: (%s %s) - %s\n", user.FullName, user.Login, err.Error())
}
