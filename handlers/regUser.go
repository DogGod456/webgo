package handlers

import (
	_ "crypto/rand"
	_ "database/sql"
	_ "encoding/hex"
	_ "github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

func RegUser(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/regUser.html")
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("user_name")
		password := r.FormValue("user_password")
		email := r.FormValue("user_email")

		hashedPassword := HashPassword(password)

		_, err := db.Exec("INSERT INTO users (user_name, user_password, user_email, hashed_password) VALUES ($1, $2, $3, $4)", username, password, email, hashedPassword)
		if err != nil {
			http.Error(w, "Ошибка регистрации. Возможно, имя пользователя или почта уже используются.", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, nil)
}
