package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

func createConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	// подключаемся к бд
	return sql.Open("postgres", connStr)
}

// getUsersFromDB получаем пользователей из базы
func getUsersFromDB() ([]*User, error) {
	users := []*User{}
	db, err := createConnection()
	if err != nil {
		return users, err
	}
	defer db.Close()
	// получаем данные агентов
	rows, err := db.Query("select id,login,first_name,last_name from users")
	if err != nil {
		return users, err
	}
	defer rows.Close()
	// считываем данные
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Login, &user.FirstName, &user.LastName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, &user)
	}

	return users, nil
}

// createUser создаем пользователя в бд
func createUser(user *UserLDAP, wg *sync.WaitGroup) error {
	defer wg.Done()
	// подключаемся к бд
	userType := "Агент"
	db, err := createConnection()
	if err != nil {
		user.errorLog(err, userType)
		return err
	}
	defer db.Close()
	// создаем пользователя
	sqlQuery := "insert into users (login,pw,title,first_name,last_name,valid_id,create_time," +
		"create_by,change_time,change_by) values ($1, $2, 'Mr/Ms', $3, $4, 1, $5, 1, $5, 1)"
	_, err = db.Exec(sqlQuery, user.Login, user.Password, user.FirstName, user.LastName, time.Now())
	if err != nil {
		user.errorLog(err, userType)
		return err
	}

	user.successLog(userType)
	return nil
}

// getCustomerUsersFromDB получаем клиентов из базы
func getCustomerUsersFromDB() ([]*CustomerUser, error) {
	users := []*CustomerUser{}
	db, err := createConnection()
	if err != nil {
		return users, err
	}
	defer db.Close()
	// получаем данные агентов
	rows, err := db.Query("select id,login,email,first_name,last_name from customer_user")
	if err != nil {
		return users, err
	}
	defer rows.Close()
	// считываем данные
	for rows.Next() {
		user := CustomerUser{}
		err := rows.Scan(&user.ID, &user.Login, &user.Email, &user.FirstName, &user.LastName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, &user)
	}

	return users, nil
}

// createCustomerUser создаем клиента в бд
func createCustomerUser(user *UserLDAP) error {
	// подключаемся к бд
	db, err := createConnection()
	if err != nil {
		user.errorLog(err, "")
		return err
	}
	defer db.Close()
	// создаем пользователя
	login := getCustomerUserLogin(user.Login)
	sqlQuery := "insert into customer_user (login,email,customer_id,first_name,last_name,phone,valid_id,create_time," +
		"create_by,change_time,change_by) values ($1, $2, $3, $4, $5, $6, 1, $7, 1, $7, 1)"
	_, err = db.Exec(sqlQuery, login, user.Login, user.Company, user.FirstName, user.LastName, user.Phone, time.Now())
	if err != nil {
		user.errorLog(err, "")
		return err
	}

	user.successLog("")
	return nil
}

func updateCustomerUser(user *UserLDAP) error {
	// подключаемся к бд
	db, err := createConnection()
	userType := "Пользователь.Обновление"

	if err != nil {
		user.errorLog(err, userType)
		return err
	}
	defer db.Close()
	// создаем пользователя
	sqlQuery := "update customer_user set customer_id = $1 where email = $2"
	_, err = db.Exec(sqlQuery, user.Company, user.Login)
	if err != nil {
		user.errorLog(err, userType)
		return err
	}

	user.successLog(userType)
	return nil
}

func getCustomerUserLogin(mail string) string {
	s := strings.Split(mail, "@")
	if len(s) == 0 {
		return mail
	}

	return s[0]
}
