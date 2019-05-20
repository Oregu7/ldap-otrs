package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/ldap.v3"
)

// User данные пользователя
type User struct {
	FullName, LastName, Username, Company, Password, Phone string
}

func main() {
	// загружаем переменные окружения
	if godotenv.Load() != nil {
		log.Fatal("Error loading .env file")
	}
	// подключаемся к серверу ldap
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	l, err := ldap.DialTLS("tcp", os.Getenv("LDAP_HOST"), tlsConfig)
	// подключаем пользователя
	err = l.Bind(os.Getenv("LDAP_LOGIN"), os.Getenv("LDAP_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}
	// ищем данные
	result, err := l.Search(&ldap.SearchRequest{
		BaseDN:       os.Getenv("LDAP_BASEDN"),
		Scope:        2,
		DerefAliases: 0,
		SizeLimit:    300,
		Filter:       "(userPrincipalName=*)",
		Attributes:   []string{"mail", "sn", "giveName", "userPrincipalName", "displayName", "company", "password", "telephoneNumber"},
		Controls:     nil,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[Количество записей: %d]\n\n", len(result.Entries))
	for indx, entry := range result.Entries {
		// пропускаем без фамилии
		if len(entry.GetAttributeValue("sn")) == 0 {
			continue
		}

		user := User{
			entry.GetAttributeValue("displayName"),
			entry.GetAttributeValue("sn"),
			entry.GetAttributeValue("userPrincipalName"),
			entry.GetAttributeValue("company"),
			entry.GetAttributeValue("password"),
			entry.GetAttributeValue("telephoneNumber"),
		}
		fmt.Println(indx, ") ", user)
	}

	fmt.Scanln()
}
