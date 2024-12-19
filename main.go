package main

import (
	"database/sql"
	"log"
	"net/http"
	"webgo/handlers"
)

func handleRequest() {
	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/register", handlers.RegUser)      // Страница регистрации
	http.HandleFunc("/login", handlers.LoginUser)       // Страница входа в аккаунт
	http.HandleFunc("/homeUser", handlers.HomeUser)     // Домашняя страница пользователя
	http.HandleFunc("/adminPage", handlers.AdminPage)   // Админ панель
	http.HandleFunc("/logout", handlers.Logout)         // Выход
	http.HandleFunc("/notes", handlers.HandleNotes)     // Общие обработчики заметок
	http.HandleFunc("/notes/", handlers.HandleNoteByID) // Частные обработчики заметок
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	var err error
	connStr := "user=postgres password=abobavgo dbname=webgodb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}

	handlers.SetDB(db)
	handleRequest()
}
