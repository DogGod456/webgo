package handlers

import (
	"database/sql"
	"errors"
	"log"
)

type Note struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Color   string `json:"color"`
}

type User struct {
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashedPassword"`
	Notes          []Note `json:"notes"`
}

func UserInformationFromDB(username string) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, user_name, hashed_password FROM users WHERE user_name = $1", username).Scan(&user.UserID, &user.Username, &user.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No user with username " + username + " found")
			return User{}, nil
		}
		log.Println("Error fetching user: ", err)
		return User{}, err
	}
	return user, nil
}
