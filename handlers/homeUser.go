package handlers

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/homeUser.html"))

func HomeUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	userID := session.Values["user_id"] // Извлекаем ID пользователя из сессии
	log.Println(userID, store, session)
	var user User

	err := db.QueryRow("SELECT id, user_name FROM users WHERE user_name = $1", userID).Scan(&user.UserID, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка получения данных1", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, user_id, title, content FROM notes WHERE user_id = $1", user.UserID)
	if err != nil {
		http.Error(w, "Ошибка получения данных2", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content); err != nil {
			http.Error(w, "Ошибка чтения данных заметки", http.StatusInternalServerError)
			return
		}
		note.UserID = user.UserID             // Устанавливаем UserID для заметки
		user.Notes = append(user.Notes, note) // Добавляем заметку в список заметок пользователя
		log.Println(user.UserID, user.Notes[0].ID, note.Title, note.Content)
	}

	err = templates.ExecuteTemplate(w, "homeUser.html", user)
	/*tmpl, err := template.ParseFiles("templates/homeUser.html")
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(user)
	tmpl.Execute(w, user)*/
}
