package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// User данные пользователя
type User struct {
	ID                         int
	Login, FirstName, LastName string
}

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
	db, err := createConnection()
	if err != nil {
		user.errorLog(err)
		return err
	}
	defer db.Close()

	sqlQuery := "insert into users (login,pw,title,first_name,last_name,valid_id,create_time," +
		"create_by,change_time,change_by) values ($1, $2, 'Mr/Ms', $3, $4, 1, $5, 1, $5, 1)"
	_, err = db.Exec(sqlQuery, user.Login, user.Password, user.FirstName, user.LastName, time.Now())
	if err != nil {
		user.errorLog(err)
		return err
	}

	user.successLog()
	return nil
}
