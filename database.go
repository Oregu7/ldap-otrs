package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// User данные пользователя
type User struct {
	FullName, LastName, Username, Company, Mail, Phone string
}

func getUsersFromDB() ([]*User, error) {
	connStr := "user=postgres password=mypass dbname=productdb sslmode=disable"
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
	users := []User{}

	for rows.Next() {
		user := user{}
		err := rows.Scan(&p.id, &p.model, &p.company, &p.price)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, user)
	}

}
