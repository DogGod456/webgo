package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/confirmLogout.html")
	if err != nil {
		log.Println("Ошибка загрузки шаблона:", err)
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)

}
