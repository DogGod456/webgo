package handlers

import (
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/homeUser.html"))

func HomeUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	userID := session.Values["user_id"] // Извлекаем ID пользователя из сессии

	if userID == nil {
		log.Println("Отсутствует сессия")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return // Возвращаем ошибку, если произошла ошибка
	}

	log.Println("asdasd", userID, store, session)
	var user User
	user.Username = session.Values["username"].(string)

	rows, err := db.Query("SELECT id, user_id, title, content FROM notes WHERE user_id = $1", userID)
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
