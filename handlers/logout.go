package handlers

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/confirmLogout.html")
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	var user User
	username := r.URL.Query().Get("username")
	user.Username = username
	user.HashedPassword = r.URL.Query().Get("&")

	if r.Method == http.MethodPost {
		var hashedPassword string

		err := db.QueryRow("SELECT hashed_password FROM users WHERE user_name = $1", username).Scan(&hashedPassword)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(w, r, "/register?error=not_registered", http.StatusSeeOther)
				return
			}
			http.Error(w, "Ошибка входа", http.StatusInternalServerError)
			log.Println("Database error:", err)
			return
		}

		if r.FormValue("confirm") == "yes" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else if r.FormValue("confirm") == "no" {
			http.Redirect(w, r, "/homeUser?username="+username+"&"+hashedPassword, http.StatusSeeOther)
			return
		}
	}

	tmpl.Execute(w, user)
}
