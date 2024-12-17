package handlers

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/loginUser.html")
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var hashedPassword string

		err = db.QueryRow("SELECT hashed_password FROM users WHERE user_name = $1", username).Scan(&hashedPassword)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(w, r, "/register?error=not_registered", http.StatusSeeOther)
				return
			}
			http.Error(w, "Ошибка входа", http.StatusInternalServerError)
			return
		}

		if CheckPasswordHash(password, hashedPassword) { // Предполагается наличие функции CheckPasswordHash в вашем коде.
			session, _ := store.Get(r, "session-name")
			session.Values["user_id"] = username // Сохраняем ID пользователя в сессии.
			session.Save(r, w)

			http.Redirect(w, r, "/homeUser", http.StatusSeeOther)
			return
		} else {
			http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
			return
		}
	}

	tmpl.Execute(w, nil)
}
