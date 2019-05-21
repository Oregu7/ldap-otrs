package main

import (
	"crypto/tls"
	"os"

	ldap "gopkg.in/ldap.v3"
)

// UserLDAP данные пользователя
type UserLDAP struct {
	FullName, LastName, Username, Company, Mail, Phone string
}

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
		Attributes:   []string{"mail", "sn", "giveName", "userPrincipalName", "displayName", "company", "mail", "telephoneNumber"},
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
			entry.GetAttributeValue("userPrincipalName"),
			entry.GetAttributeValue("company"),
			entry.GetAttributeValue("mail"),
			entry.GetAttributeValue("telephoneNumber"),
		}
		users = append(users, user)
	}

	return users, nil
}
