package handlers

import (
	_ "database/sql"
	"html/template"
	"log"
	"net/http"
)

func OpenNote(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	hashedPassword := r.URL.Query().Get("&")

	// Получаем список заметок из базы данных для данного пользователя
	rows, err := db.Query("SELECT id, title FROM notes WHERE user_id = (SELECT id FROM users WHERE user_name = $1 and hashed_password = $2)", username, hashedPassword)
	if err != nil {
		http.Error(w, "Ошибка получения списка заметок или доступа", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.UserID, &note.Title); err != nil {
			log.Println("Ошибка сканирования заметки:", err)
			continue
		}
		notes = append(notes, note)
	}

	tmpl, err := template.ParseFiles("templates/openNote.html")
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{
		"Username": username,
		"Notes":    notes,
	})
}
