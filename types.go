package main

import "fmt"

// UserHashed хэшируемы пользователь
type UserHashed interface {
	getHashMapKey() string
	getPropsToken() string
}

// UserLDAP данные пользователя
type UserLDAP struct {
	FullName, LastName, FirstName, Login string
	Company, Mail, Phone                 string
	Password                             string
}

// User данные пользователя
type User struct {
	ID                         int
	Login, FirstName, LastName string
}

// CustomerUser данные клиента
type CustomerUser struct {
	ID                         int
	Login, Email, CustomerID   string
	FirstName, LastName, Phone string
}

func (user *UserLDAP) successLog(userType string) {
	if userType == "" {
		userType = "Пользователь"
	}
	fmt.Printf("[OK] %s: (%s %s)\n", userType, user.FullName, user.Login)
}

func (user *UserLDAP) errorLog(err error, userType string) {
	if userType == "" {
		userType = "Пользователь"
	}
	fmt.Printf("[ERROR] %s: (%s %s) - %s\n", userType, user.FullName, user.Login, err.Error())
}

// реализуем интерфейс UserHashed

func (user *UserLDAP) getPropsToken() string {
	token := fmt.Sprintf("%s", user.Company)
	return token
}

func (user *UserLDAP) getHashMapKey() string {
	return user.Login
}

func (user *CustomerUser) getPropsToken() string {
	token := fmt.Sprintf("%s", user.CustomerID)
	return token
}

func (user *CustomerUser) getHashMapKey() string {
	return user.Email
}
