package handlers

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
)

func HomeUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	log.Println("homeuser", username)
	// Получаем данные о пользователе из базы данных
	var user User
	err := db.QueryRow("SELECT user_name, hashed_password FROM users WHERE user_name = $1", username).Scan(&user.Username, &user.HashedPassword)
	log.Println(err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	// Отображаем страницу с данными пользователя
	tmpl, err := template.ParseFiles("templates/homeUser.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl.Execute(w, user)
}
