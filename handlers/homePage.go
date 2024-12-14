package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/homePage.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}
