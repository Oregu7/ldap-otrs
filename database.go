package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// User данные пользователя
type User struct {
	ID                         int
	Login, FirstName, LastName string
}

// getUsersFromDB получаем пользователей из базы
func getUsersFromDB() ([]*User, error) {
	users := []*User{}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	// подключаемся к бд
	db, err := sql.Open("postgres", connStr)
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
