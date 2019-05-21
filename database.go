package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// User данные пользователя
type User struct {
	ID                                          int
	Login, Password, Title, FirstName, LastName string
}

// getUsersFromDB получаем пользователей из базы
func getUsersFromDB() ([]*User, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	users := []*User{}

	for rows.Next() {
		user := &User{}
		err := rows.Scan(user.ID, user.Login, user.Password, user.Title, user.FirstName, user.LastName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}
