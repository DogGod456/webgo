package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	delete(session.Values, "user_id") // Удаляем ID пользователя из сессии
	session.Save(r, w)

	tmpl, err := template.ParseFiles("templates/confirmLogout.html")
	if err != nil {
		log.Println("Ошибка загрузки шаблона:", err)
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
