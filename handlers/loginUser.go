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

		hashedPassword := HashPassword(password)

		err = db.QueryRow("SELECT hashed_password FROM users WHERE user_name = $1", username).Scan(&hashedPassword)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(w, r, "/register?error=not_registered", http.StatusSeeOther)
				return
			}
			http.Error(w, "Ошибка входа", http.StatusInternalServerError)
			panic(err)
			return
		}
		log.Println(CheckPasswordHash(password, hashedPassword))
		if CheckPasswordHash(password, hashedPassword) {
			http.Redirect(w, r, "/register?error=invalid_credentials", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/homeUser?username="+username+"&"+hashedPassword, http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, nil)
}
